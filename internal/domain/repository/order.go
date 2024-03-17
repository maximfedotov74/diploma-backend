package repository

import (
	"context"
	"errors"
	"fmt"

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

func (r *OrderRepository) Create(ctx context.Context, input model.CreateOrderInput, userId int) (*model.CreateOrderResponse, fall.Error) {

	var ex fall.Error = nil

	tx, err := r.db.Begin(ctx)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
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
		return nil, ex
	}

	for _, item := range input.CartItems {
		ex = r.createOrderModel(ctx, tx, orderId, item)
		if ex != nil {
			return nil, ex
		}
	}

	ex = r.AddDeliveryPoint(ctx, tx, orderId, input.DeliveryPointId)
	if ex != nil {
		return nil, ex
	}

	link, ex := r.AddActivationLink(ctx, tx, orderId)
	if ex != nil {
		return nil, ex
	}

	if len(input.CartItems) > 0 {
		var cartIds []int
		for _, item := range input.CartItems {
			cartIds = append(cartIds, item.CartItemId)
		}
		ex = r.wishRepository.RemoveSeveralItems(ctx, tx, cartIds)
		if ex != nil {
			return nil, ex
		}
	}

	return &model.CreateOrderResponse{Link: *link, Id: orderId, Total: input.TotalPrice}, nil
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

func (or *OrderRepository) ActivateOrder(link string) fall.Error {
	ctx := context.Background()
	query := `SELECT order_id FROM order_activation WHERE link = $1 AND end_time > CURRENT_TIMESTAMP;`

	row := or.db.QueryRow(ctx, query, link)

	var orderId string

	err := row.Scan(&orderId)

	if err != nil {
		return fall.ServerError(msg.OrderErrorWhenActivate)
	}

	query = "UPDATE public.order SET is_activated = TRUE WHERE order_id = $1;"

	_, err = or.db.Exec(ctx, query, orderId)
	if err != nil {
		return fall.ServerError(msg.OrderErrorWhenActivate)
	}

	return nil

}

func (or *OrderRepository) GetUserOrders(ctx context.Context, userId int) ([]*model.Order, fall.Error) {

	query := `
	SELECT o.order_id as o_id, o.created_at as o_created_at, o.updated_at as o_updated_at,
	o.delivery_date as o_delivery_date, o.is_activated as o_is_activated,
	o.order_status as o_order_status, o.order_payment_method as o_payment_method,
	o.conditions as o_conditions, o.products_price as o_products_price,
	o.total_price as o_total_price, o.total_discount as o_total_discount,o.promo_discount as o_promo_discount,
	o.delivery_price as o_delivery_price,o.recipient_firstname as o_recipient_firstname, 
	o.recipient_lastname as o_recipient_lastname,o.recipient_phone as o_recipient_phone, u.user_id as u_id, u.email as u_email,
	om.order_model_id as om_id, om.quantity as om_quantity,
	om.price as om_price, om.discount as om_discount,
	ms.model_size_id as ms_id, ms.product_model_id as ms_product_model_id, ms.size_id as ms_size_id, ms.literal_size as ms_literal_size,
	sz.size_value as ms_size_value,  ms.in_stock as ms_in_stock,
	pm.main_image_path as pm_main_image_path,
	p.product_id as p_id, p.title as p_title, pm.product_model_id as model_id, pm.slug as pm_slug, pm.article as pm_article,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug,
	dp.delivery_point_id as dp_id, dp.title as dp_title, dp.city as dp_city,
	dp.address as dp_address, dp.coords as db_coords, dp.with_fitting as dp_with_fitting,
	dp.work_schedule as dp_work_schedule, dp.info as dp_info
	FROM public.order as o
	INNER JOIN public.user as u ON o.user_id = u.user_id
	INNER JOIN order_model as om ON o.order_id = om.order_id
	INNER JOIN model_sizes as ms ON om.model_size_id = ms.model_size_id
	INNER JOIN product_model as pm ON ms.product_model_id = pm.product_model_id
	INNER JOIN sizes as sz ON ms.size_id = sz.size_id
	INNER JOIN product as p ON p.product_id = pm.product_id
	INNER JOIN category as c on p.category_id = c.category_id
	INNER JOIN brand as b on p.brand_id = b.brand_id
	INNER JOIN order_delivery_point as odp ON o.order_id = odp.order_id
	INNER JOIN delivery_point as dp ON odp.delivery_point_id = dp.delivery_point_id
	WHERE u.user_id = $1
	ORDER BY o.created_at DESC;
	`

	rows, err := or.db.Query(ctx, query, userId)

	if err != nil {

		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	ordersMap := make(map[string]*model.Order)
	var ordersOrder []string
	for rows.Next() {
		o := model.Order{}
		m := model.OrderModel{}

		err := rows.Scan(&o.Id, &o.CreatedAt, &o.UpdatedAt, &o.DeliveryDate, &o.IsActivated, &o.Status, &o.PaymentMethod, &o.Conditions,
			&o.ProductsPrice, &o.TotalPrice, &o.TotalDiscount, &o.PromoDiscount, &o.DeliveryPrice, &o.User.FirstName, &o.User.LastName,
			&o.User.Phone, &o.User.Id, &o.User.Email, &m.OrderModelId, &m.Quantity, &m.Price, &m.Discount, &m.Size.ModelId, &m.Size.ModelId, &m.Size.SizeId, &m.Size.Literal, &m.Size.Value, &m.Size.InStock, &m.MainImagePath, &m.Product.ProductId, &m.Product.Title, &m.ModelId, &m.Slug, &m.Article,
			&m.Product.Category.Id, &m.Product.Category.Title, &m.Product.Category.Slug,
			&m.Product.Brand.Id, &m.Product.Brand.Title, &m.Product.Brand.Slug, &o.DeliveryPoint.Id, &o.DeliveryPoint.Title,
			&o.DeliveryPoint.City, &o.DeliveryPoint.Address, &o.DeliveryPoint.Coords, &o.DeliveryPoint.WithFitting, &o.DeliveryPoint.WorkSchedule,
			&o.DeliveryPoint.Info,
		)
		if err != nil {

			return nil, fall.ServerError(err.Error())
		}

		current, ok := ordersMap[o.Id]
		if !ok {
			o.Models = append(o.Models, m)
			ordersMap[o.Id] = &o
			ordersOrder = append(ordersOrder, o.Id)
		} else {
			current.Models = append(current.Models, m)
			ordersMap[current.Id] = current
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	orders := make([]*model.Order, 0, len(ordersMap))

	for _, id := range ordersOrder {
		o := ordersMap[id]
		orders = append(orders, o)
	}

	return orders, nil
}

func (or *OrderRepository) GetOrder(ctx context.Context, id string) (*model.Order, fall.Error) {

	query := `
	SELECT o.order_id as o_id, o.created_at as o_created_at, o.updated_at as o_updated_at,
	o.delivery_date as o_delivery_date, o.is_activated as o_is_activated,
	o.order_status as o_order_status, o.order_payment_method as o_payment_method,
	o.conditions as o_conditions, o.products_price as o_products_price,
	o.total_price as o_total_price, o.total_discount as o_total_discount,o.promo_discount as o_promo_discount,
	o.delivery_price as o_delivery_price,o.recipient_firstname as o_recipient_firstname, 
	o.recipient_lastname as o_recipient_lastname,o.recipient_phone as o_recipient_phone, u.user_id as u_id, u.email as u_email,
	om.order_model_id as om_id, om.quantity as om_quantity,
	om.price as om_price, om.discount as om_discount,
	ms.model_size_id as ms_id, ms.product_model_id as ms_product_model_id, ms.size_id as ms_size_id, ms.literal_size as ms_literal_size,
	sz.size_value as ms_size_value, ms.in_stock as ms_in_stock,
	pm.main_image_path as pm_main_image_path,
	p.product_id as p_id, p.title as p_title, pm.slug as pm_slug, pm.article as pm_article,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug,
	dp.delivery_point_id as dp_id, dp.title as dp_title, dp.city as dp_city,
	dp.address as dp_address, dp.coords as db_coords, dp.with_fitting as dp_with_fitting,
	dp.work_schedule as dp_work_schedule, dp.info as dp_info
	FROM public.order as o
	INNER JOIN public.user as u ON o.user_id = u.user_id
	INNER JOIN order_model as om ON o.order_id = om.order_id
	INNER JOIN model_sizes as ms ON om.model_size_id = ms.model_size_id
	INNER JOIN product_model as pm ON ms.product_model_id = pm.product_model_id
	INNER JOIN sizes as sz ON ms.size_id = sz.size_id
	INNER JOIN product as p ON p.product_id = pm.product_id
	INNER JOIN category as c on p.category_id = c.category_id
	INNER JOIN brand as b on p.brand_id = b.brand_id
	INNER JOIN order_delivery_point as odp ON o.order_id = odp.order_id
	INNER JOIN delivery_point as dp ON odp.delivery_point_id = dp.delivery_point_id
	WHERE o.order_id = $1;
	`

	rows, err := or.db.Query(ctx, query, id)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	o := model.Order{}

	founded := false

	for rows.Next() {
		m := model.OrderModel{}
		err := rows.Scan(&o.Id, &o.CreatedAt, &o.UpdatedAt, &o.DeliveryDate, &o.IsActivated, &o.Status, &o.PaymentMethod, &o.Conditions,
			&o.ProductsPrice, &o.TotalPrice, &o.TotalDiscount, &o.PromoDiscount, &o.DeliveryPrice, &o.User.FirstName, &o.User.LastName,
			&o.User.Phone, &o.User.Id, &o.User.Email, &m.OrderModelId, &m.Quantity, &m.Price, &m.Discount, &m.Size.SizeModelId, &m.Size.ModelId, &m.Size.SizeId, &m.Size.Literal, &m.Size.Value, &m.Size.InStock, &m.MainImagePath, &m.Product.ProductId, &m.Product.Title, &m.Slug, &m.Article,
			&m.Product.Category.Id, &m.Product.Category.Title, &m.Product.Category.Slug,
			&m.Product.Brand.Id, &m.Product.Brand.Title, &m.Product.Brand.Slug, &o.DeliveryPoint.Id, &o.DeliveryPoint.Title,
			&o.DeliveryPoint.City, &o.DeliveryPoint.Address, &o.DeliveryPoint.Coords, &o.DeliveryPoint.WithFitting, &o.DeliveryPoint.WorkSchedule,
			&o.DeliveryPoint.Info,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		o.Models = append(o.Models, m)

		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.OrderNotFound, fall.STATUS_NOT_FOUND)
	}

	return &o, nil
}

func (or *OrderRepository) GetAdminOrders(ctx context.Context, page int, fromDate *string, toDate *string) (*model.AllOrdersResponse, fall.Error) {

	limit := 24

	offset := page*limit - limit

	whereFilter := ""

	if fromDate != nil && toDate == nil {
		whereFilter = fmt.Sprintf("WHERE o.created_at >= '%s'", *fromDate)
	}

	if fromDate == nil && toDate != nil {
		whereFilter = fmt.Sprintf("WHERE o.created_at <= '%s'", *toDate)
	}

	if fromDate != nil && toDate != nil {
		whereFilter = fmt.Sprintf("WHERE o.created_at >= '%s' AND o.created_at <= '%s'", *fromDate, *toDate)

	}

	query := fmt.Sprintf(`
	SELECT DISTINCT o.order_id as o_id, o.created_at,
	(SELECT count(distinct o.order_id)
	FROM public.order as o
	INNER JOIN public.user as u ON o.user_id = u.user_id
	INNER JOIN order_model as om ON o.order_id = om.order_id
	INNER JOIN model_sizes as ms ON om.model_size_id = ms.model_size_id
	INNER JOIN product_model as pm ON ms.product_model_id = pm.product_model_id
	INNER JOIN sizes as sz ON ms.size_id = sz.size_id
	INNER JOIN product as p ON p.product_id = pm.product_id
	INNER JOIN category as c on p.category_id = c.category_id
	INNER JOIN brand as b on p.brand_id = b.brand_id
	INNER JOIN order_delivery_point as odp ON o.order_id = odp.order_id
	INNER JOIN delivery_point as dp ON odp.delivery_point_id = dp.delivery_point_id
	) as total
	FROM public.order as o
	INNER JOIN public.user as u ON o.user_id = u.user_id
	INNER JOIN order_model as om ON o.order_id = om.order_id
	INNER JOIN model_sizes as ms ON om.model_size_id = ms.model_size_id
	INNER JOIN product_model as pm ON ms.product_model_id = pm.product_model_id
	INNER JOIN sizes as sz ON ms.size_id = sz.size_id
	INNER JOIN product as p ON p.product_id = pm.product_id
	INNER JOIN category as c on p.category_id = c.category_id
	INNER JOIN brand as b on p.brand_id = b.brand_id
	INNER JOIN order_delivery_point as odp ON o.order_id = odp.order_id
	INNER JOIN delivery_point as dp ON odp.delivery_point_id = dp.delivery_point_id
	%s
	ORDER BY o.created_at DESC
	LIMIT $1 OFFSET $2
	;`, whereFilter)

	rows, err := or.db.Query(ctx, query, limit, offset)

	if err != nil {

		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var total int
	var ordersOrder []string

	for rows.Next() {
		var orderId string
		err := rows.Scan(&orderId, nil, &total)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		ordersOrder = append(ordersOrder, orderId)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	query = `
	SELECT o.order_id as o_id, o.created_at as o_created_at, o.updated_at as o_updated_at,
	o.delivery_date as o_delivery_date, o.is_activated as o_is_activated,
	o.order_status as o_order_status, o.order_payment_method as o_payment_method,
	o.conditions as o_conditions, o.products_price as o_products_price,
	o.total_price as o_total_price, o.total_discount as o_total_discount,o.promo_discount as o_promo_discount,
	o.delivery_price as o_delivery_price,o.recipient_firstname as o_recipient_firstname, 
	o.recipient_lastname as o_recipient_lastname,o.recipient_phone as o_recipient_phone, u.user_id as u_id, u.email as u_email,
	om.order_model_id as om_id, om.quantity as om_quantity,
	om.price as om_price, om.discount as om_discount,
	ms.model_size_id as ms_id, ms.product_model_id as ms_product_model_id, ms.size_id as ms_size_id, ms.literal_size as ms_literal_size,
	sz.size_value as ms_size_value, ms.in_stock as ms_in_stock,
	pm.main_image_path as pm_main_image_path,
	p.product_id as p_id, p.title as p_title, pm.product_model_id as model_id, pm.slug as pm_slug, pm.article as pm_article,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug,
	dp.delivery_point_id as dp_id, dp.title as dp_title, dp.city as dp_city,
	dp.address as dp_address, dp.coords as db_coords, dp.with_fitting as dp_with_fitting,
	dp.work_schedule as dp_work_schedule, dp.info as dp_info
	FROM public.order as o
	INNER JOIN public.user as u ON o.user_id = u.user_id
	INNER JOIN order_model as om ON o.order_id = om.order_id
	INNER JOIN model_sizes as ms ON om.model_size_id = ms.model_size_id
	INNER JOIN product_model as pm ON ms.product_model_id = pm.product_model_id
	INNER JOIN sizes as sz ON ms.size_id = sz.size_id
	INNER JOIN product as p ON p.product_id = pm.product_id
	INNER JOIN category as c on p.category_id = c.category_id
	INNER JOIN brand as b on p.brand_id = b.brand_id
	INNER JOIN order_delivery_point as odp ON o.order_id = odp.order_id
	INNER JOIN delivery_point as dp ON odp.delivery_point_id = dp.delivery_point_id
	WHERE o.order_id = ANY ($1);
	`

	rows, err = or.db.Query(ctx, query, ordersOrder)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	ordersMap := make(map[string]*model.Order)
	for rows.Next() {
		o := model.Order{}
		m := model.OrderModel{}

		err := rows.Scan(&o.Id, &o.CreatedAt, &o.UpdatedAt, &o.DeliveryDate, &o.IsActivated, &o.Status, &o.PaymentMethod, &o.Conditions,
			&o.ProductsPrice, &o.TotalPrice, &o.TotalDiscount, &o.PromoDiscount, &o.DeliveryPrice, &o.User.FirstName, &o.User.LastName,
			&o.User.Phone, &o.User.Id, &o.User.Email, &m.OrderModelId, &m.Quantity, &m.Price, &m.Discount, &m.Size.SizeModelId, &m.Size.ModelId, &m.Size.SizeId, &m.Size.Literal, &m.Size.Value, &m.Size.InStock, &m.MainImagePath, &m.Product.ProductId, &m.Product.Title, &m.ModelId, &m.Slug, &m.Article,
			&m.Product.Category.Id, &m.Product.Category.Title, &m.Product.Category.Slug,
			&m.Product.Brand.Id, &m.Product.Brand.Title, &m.Product.Brand.Slug, &o.DeliveryPoint.Id, &o.DeliveryPoint.Title,
			&o.DeliveryPoint.City, &o.DeliveryPoint.Address, &o.DeliveryPoint.Coords, &o.DeliveryPoint.WithFitting, &o.DeliveryPoint.WorkSchedule,
			&o.DeliveryPoint.Info,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		current, ok := ordersMap[o.Id]
		if !ok {
			o.Models = append(o.Models, m)
			ordersMap[o.Id] = &o
		} else {
			current.Models = append(current.Models, m)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	orders := make([]*model.Order, 0, len(ordersMap))

	for _, id := range ordersOrder {
		o := ordersMap[id]
		orders = append(orders, o)
	}

	return &model.AllOrdersResponse{
		Orders: orders, Total: total,
	}, nil
}

func (or *OrderRepository) SendNewActivationLink(ctx context.Context, orderId string) (*string, fall.Error) {

	var ex fall.Error = nil

	tx, err := or.db.Begin(ctx)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}
	defer func() {
		if ex != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := "SELECT is_activated FROM public.order WHERE order_id = $1;"

	row := tx.QueryRow(ctx, query, orderId)

	var isActivated bool

	err = row.Scan(&isActivated)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	if isActivated {
		ex = fall.NewErr(msg.OrderAlreadyActivated, fall.STATUS_BAD_REQUEST)
		return nil, ex
	}
	link, ex := or.AddActivationLink(ctx, tx, orderId)
	if ex != nil {
		return nil, ex
	}
	return link, nil
}

func (r *OrderRepository) CancelOrder(ctx context.Context, orderId string, userId int) fall.Error {
	q := "UPDATE public.order SET order_status = 'canceled' WHERE order_id = $1 AND user_id = $2;"

	_, err := r.db.Exec(ctx, q, orderId, userId)

	if err != nil {
		return fall.ServerError(msg.OrderErrorWhenCancel)
	}

	return nil
}

func (r *OrderRepository) ChangeStatus(ctx context.Context, orderId string, status model.OrderStatusEnum) fall.Error {
	q := "UPDATE public.order SET order_status = $1 WHERE order_id = $2;"

	_, err := r.db.Exec(ctx, q, status, orderId)

	if err != nil {
		return fall.ServerError(msg.OrderErrorWhenChangeStatus)
	}

	return nil
}
