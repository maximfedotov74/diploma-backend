package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type orderProductRepository interface {
	ReturnQuantityInStock(ctx context.Context, modelSizeId int, quantity int, tx db.Transaction) fall.Error
	ReduceQuantityInStock(ctx context.Context, modelSizeId int, quantity int, tx db.Transaction) fall.Error
}

type orderWishRepository interface {
	RemoveSeveralItems(ctx context.Context, tx db.Transaction, cartIds []int) fall.Error
}

type OrderRepository struct {
	db                db.PostgresClient
	wishRepository    orderWishRepository
	productRepository orderProductRepository
}

func NewOrderRepository(db db.PostgresClient, wishRepository orderWishRepository,
	productRepository orderProductRepository) *OrderRepository {
	return &OrderRepository{db: db, wishRepository: wishRepository, productRepository: productRepository}
}

func (r *OrderRepository) Create(ctx context.Context, input model.CreateOrderInput, userId int) fall.Error {

	var ex fall.Error

	tx, err := r.db.Begin(ctx)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return ex
	}

	defer func() {
		if ex != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	var status model.OrderStatusEnum = model.WaitingForActivation

	if input.PaymentMethod == model.Online {
		status = model.WaitingForPayment
	}

	query := `INSERT INTO public.order (order_payment_method,conditions,products_price,total_price,total_discount,delivery_price,recipient_firstname,recipient_lastname,recipient_phone,user_id, order_status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	RETURNING order_id;`

	row := tx.QueryRow(ctx, query, input.PaymentMethod, input.Conditions, input.ProductsPrice, input.TotalPrice, input.TotalDiscount, input.DeliveryPrice, input.RecipientFirstname, input.RecipientLastname, input.RecipientPhone, userId, status)

	var orderId string

	err = row.Scan(&orderId)
	if err != nil {
		ex = fall.ServerError(err.Error())
		return ex
	}

	for _, item := range input.CartItems {
		ex = r.createOrderModel(ctx, tx, orderId, item)
		if ex != nil {
			return ex
		}
	}

	ex = r.AddDeliveryPoint(ctx, tx, orderId, input.DeliveryPointId)
	if ex != nil {
		return ex
	}

	link, ex := r.AddActivationLink(ctx, tx, orderId)
	if ex != nil {
		return ex
	}

	if len(input.CartItems) > 0 {
		var cartIds []int
		for _, item := range input.CartItems {
			cartIds = append(cartIds, item.CartItemId)
		}
		ex = r.wishRepository.RemoveSeveralItems(ctx, tx, cartIds)
		if ex != nil {
			return ex
		}
	}

	log.Println(orderId, link, input.TotalPrice)

	return nil
}

func (r *OrderRepository) createOrderModel(
	ctx context.Context, tx db.Transaction,
	orderId string, item *model.CartItemModel) fall.Error {

	query := `
	INSERT INTO order_model (order_id,model_size_id,quantity,price,discount) VALUES ($1,$2,$3,$4,$5); 
	`
	_, err := tx.Exec(ctx, query, orderId, item.ModelSizeId, item.Quantity, item.Price, item.Discount)

	if err != nil {
		return fall.ServerError(msg.OrderErrorWhenAddModelsToProduct)
	}
	return nil

}

func (or *OrderRepository) AddDeliveryPoint(ctx context.Context, tx db.Transaction, orderId string, pointId int) fall.Error {
	query := `
	INSERT INTO order_delivery_point (order_id,delivery_point_id) VALUES ($1,$2);
	`
	_, err := tx.Exec(ctx, query, orderId, pointId)

	if err != nil {
		return fall.ServerError(msg.OrderErrorWhenAddDeliveryPoint)
	}

	return nil
}

func (or *OrderRepository) FindActivationLink(ctx context.Context, tx db.Transaction, orderId string) (*int, fall.Error) {
	query := "SELECT order_activation_id FROM order_activation WHERE order_id = $1;"

	row := tx.QueryRow(ctx, query, orderId)

	var activationId int

	err := row.Scan(&activationId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.OrderActivationLinkNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &activationId, nil
}

func (or *OrderRepository) UpdateActivationLink(ctx context.Context, tx db.Transaction, activationId int) (*string, fall.Error) {
	query := `UPDATE order_activation SET
	link = uuid_generate_v4(),
	end_time = CURRENT_TIMESTAMP + INTERVAL '4 hours'
	WHERE order_activation_id = $1 RETURNING link;`

	row := tx.QueryRow(ctx, query, activationId)

	var link string

	err := row.Scan(&link)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	return &link, nil
}

func (or *OrderRepository) AddActivationLink(ctx context.Context, tx db.Transaction, orderId string) (*string, fall.Error) {

	activationId, err := or.FindActivationLink(ctx, tx, orderId)

	if err != nil {
		query := `
		INSERT INTO order_activation (order_id) VALUES ($1) RETURNING link;
		`
		row := tx.QueryRow(ctx, query, orderId)
		var link string
		err := row.Scan(&link)
		if err != nil {
			return nil, fall.ServerError(msg.OrderErrorWhenAddActivationLink)
		}
		return &link, nil
	}
	link, err := or.UpdateActivationLink(ctx, tx, *activationId)
	if err != nil {
		return nil, err
	}
	return link, nil
}
