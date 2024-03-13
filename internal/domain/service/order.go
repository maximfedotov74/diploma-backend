package service

import (
	"context"
	"math"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/payment"
)

type orderRepository interface {
	Create(ctx context.Context, input model.CreateOrderInput, userId int) (*model.CreateOrderResponse, fall.Error)
	GetAdminOrders(ctx context.Context) ([]*model.Order, fall.Error)
	GetOrder(ctx context.Context, id string) (*model.Order, fall.Error)
	GetUserOrders(ctx context.Context, userId int) ([]*model.Order, fall.Error)
	CancelOrder(ctx context.Context, orderId string, userId int) fall.Error
}

type orderUserService interface {
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
}

type orderMailService interface {
	SendOrderActivationEmail(to string, subject string, code string) error
}

type orderDeliveryRepository interface {
	FindById(ctx context.Context, id int) (*model.DeliveryPoint, fall.Error)
}

type orderPaymentService interface {
	CreatePayment(orderId string, totalPrice float64) (*payment.Payment, error)
}

type orderWishService interface {
	FindModelInUserCart(ctx context.Context, modelSizeId int, userId int) (*model.CartItemModel, fall.Error)
}

type OrderService struct {
	repo           orderRepository
	wishService    orderWishService
	userService    orderUserService
	deliveryRepo   orderDeliveryRepository
	mailService    orderMailService
	paymentService orderPaymentService
}

func NewOrderService(repo orderRepository, wishService orderWishService, userService orderUserService,
	deliveryRepo orderDeliveryRepository, mailService orderMailService, paymentService orderPaymentService) *OrderService {
	return &OrderService{
		repo:           repo,
		wishService:    wishService,
		userService:    userService,
		deliveryRepo:   deliveryRepo,
		mailService:    mailService,
		paymentService: paymentService,
	}
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId string, userId int) fall.Error {
	return s.repo.CancelOrder(ctx, orderId, userId)
}

func (s *OrderService) GetAdminOrders(ctx context.Context) ([]*model.Order, fall.Error) {
	return s.repo.GetAdminOrders(ctx)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userId int) ([]*model.Order, fall.Error) {
	return s.repo.GetUserOrders(ctx, userId)
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, fall.Error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *OrderService) Create(ctx context.Context, dto model.CreateOrderDto, user *model.LocalSession) (*string, fall.Error) {

	var cartItems []*model.CartItemModel

	for _, id := range dto.ModelSizeIds {
		item, ex := s.wishService.FindModelInUserCart(ctx, id, user.UserId)
		if ex != nil {
			continue
		}
		cartItems = append(cartItems, item)
	}

	deliveryPoint, ex := s.deliveryRepo.FindById(ctx, dto.DeliveryPointId)
	if ex != nil {
		return nil, ex
	}

	var productsPrice float64 = 0
	var totalDiscount float64 = 0

	for _, item := range cartItems {
		productsPrice += (float64(item.Price) * (float64(item.Quantity)))
		if item.Discount != nil {
			totalDiscount += (float64(item.Price) / 100) * float64(*item.Discount) * float64(item.Quantity)
		}
	}

	flag := deliveryPoint.WithFitting == model.ConvertFittingToBool(dto.Conditions)
	if !flag {
		return nil, fall.NewErr(msg.OrderDeliveryPointConditionConflict, fall.STATUS_BAD_REQUEST)
	}

	var deliveryPrice float64 = 0
	if model.ConvertFittingToBool(dto.Conditions) {
		deliveryPrice = 199
	}

	totalDiscount = math.Ceil(totalDiscount)
	productsPrice = math.Ceil(productsPrice)

	totalPrice := productsPrice - totalDiscount + deliveryPrice

	input := model.CreateOrderInput{
		DeliveryPrice:      deliveryPrice,
		TotalPrice:         totalPrice,
		ProductsPrice:      productsPrice,
		TotalDiscount:      totalDiscount,
		RecipientFirstname: dto.RecipientFirstname,
		RecipientLastname:  dto.RecipientLastname,
		RecipientPhone:     dto.RecipientPhone,
		PaymentMethod:      dto.PaymentMethod,
		DeliveryPointId:    deliveryPoint.Id,
		CartItems:          cartItems,
		Conditions:         dto.Conditions,
	}

	resp, ex := s.repo.Create(ctx, input, user.UserId)
	if ex != nil {
		return nil, ex
	}

	if dto.PaymentMethod == model.Online {
		p, err := s.paymentService.CreatePayment(resp.Id, resp.Total)
		if err != nil {
			return nil, fall.ServerError("Ошибка при обработки платежа №" + resp.Id)
		}
		return &p.Confirmation.ConfirmationURL, nil
	}

	go s.mailService.SendOrderActivationEmail(user.Email, "Подтверждение оформления заказа!", resp.Link)
	return nil, nil

}
