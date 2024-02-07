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
	Create(ctx context.Context, input model.CreateOrderInput, userId int) fall.Error
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
	CreatePayment(orderId string, totalPrice int) (*payment.Payment, error)
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

func (s *OrderService) Create(ctx context.Context, dto model.CreateOrderDto, user model.LocalSession) fall.Error {

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
		return ex
	}

	productsPrice := 0
	totalDiscount := 0

	for _, item := range cartItems {
		productsPrice += (item.Price * item.Quantity)
		if item.Discount != nil {
			discountPrice := math.Round(float64(item.Price / 100 * int(*item.Discount)))
			totalDiscount += int(discountPrice) * item.Quantity
		}
	}

	flag := deliveryPoint.WithFitting == model.ConvertFittingToBool(dto.Conditions)
	if !flag {
		return fall.NewErr(msg.OrderDeliveryPointConditionConflict, fall.STATUS_BAD_REQUEST)
	}

	deliveryPrice := 0
	if model.ConvertFittingToBool(dto.Conditions) {
		deliveryPrice = 199
	}

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

	return s.repo.Create(ctx, input, user.UserId)

}
