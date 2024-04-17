package model

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type AllOrdersResponse struct {
	Orders []*Order `json:"orders"`
	Total  int      `json:"total"`
}

type CreateOrderResponse struct {
	Link  string  `json:"link" validate:"required"`
	Id    string  `json:"id" validate:"required"`
	Total float64 `json:"total" validate:"required"`
}

type ChangeOrderStatusDto struct {
	Status OrderStatusEnum `json:"status" validate:"required"`
}

type ChangeOrderDeliveryDate struct {
	Date time.Time `json:"date" validate:"required"`
}

type OrderConfirmation struct {
	PaymentUrl *string `json:"payment_url"`
}

type OrderStatusEnum string

const (
	Completed            OrderStatusEnum = "completed"
	Canceled             OrderStatusEnum = "canceled"
	OnTheWay             OrderStatusEnum = "on_the_way"
	WaitingForPayment    OrderStatusEnum = "waiting_for_payment"
	Paid                 OrderStatusEnum = "paid"
	InProcessing         OrderStatusEnum = "in_processing"
	WaitingForActivation OrderStatusEnum = "waiting_for_activation"
)

type PaymentMethodEnum string

const (
	UponReceipt PaymentMethodEnum = "upon_receipt"
	Online      PaymentMethodEnum = "online"
)

type OrderConditions string

const (
	WithFitting    OrderConditions = "with_fitting"
	WithoutFitting OrderConditions = "without_fitting"
)

type OrderUser struct {
	Id        int    `json:"id" validate:"required"`
	Email     string `json:"recipient_email" validate:"required"`
	Phone     string `json:"recipient_phone" validate:"required"`
	FirstName string `json:"recipient_firstname" validate:"required"`
	LastName  string `json:"recipient_lastname" validate:"required"`
}

type Order struct {
	Id            string            `json:"order_id" validate:"required"`
	PaymentId     *string           `json:"-"`
	User          OrderUser         `json:"user" validate:"required"`
	CreatedAt     time.Time         `json:"created_at" validate:"required"`
	UpdatedAt     time.Time         `json:"updated_at" validate:"required"`
	DeliveryDate  *time.Time        `json:"delivery_date"`
	IsActivated   bool              `json:"is_activated" validate:"required"`
	Status        OrderStatusEnum   `json:"status" validate:"required"`
	PaymentMethod PaymentMethodEnum `json:"payment_method" validate:"required"`
	Conditions    OrderConditions   `json:"conditions" validate:"required"`
	ProductsPrice float64           `json:"products_price" validate:"required"`
	TotalPrice    float64           `json:"total_price" validate:"required"`
	TotalDiscount *float64          `json:"total_discount"`
	PromoDiscount *int              `json:"promo_discount"`
	DeliveryPrice int               `json:"delivery_price" validate:"required"`
	DeliveryPoint DeliveryPoint     `json:"delivery_point" validate:"required"`
	Models        []OrderModel      `json:"models" validate:"required"`
}

type OrderModelProduct struct {
	ProductId int           `json:"product_id" validate:"required"`
	Title     string        `json:"title" validate:"required"`
	Category  CategoryModel `json:"category" validate:"required"`
	Brand     Brand         `json:"brand" validate:"required"`
}

type OrderModel struct {
	OrderModelId  int               `json:"order_model_id" validate:"required"`
	ModelId       int               `json:"model_id" validate:"required"`
	Slug          string            `json:"slug" validate:"required"`
	Article       string            `json:"article" validate:"required"`
	Quantity      int               `json:"quantity" validate:"required"`
	Price         int               `json:"price" validate:"required"`
	Discount      *byte             `json:"discount"`
	Size          ProductModelSize  `json:"size" validate:"required"`
	MainImagePath string            `json:"main_image_path" validate:"required"`
	Product       OrderModelProduct `json:"product" validate:"required"`
}

func OrderStatusEnumValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(Completed), string(Canceled), string(OnTheWay), string(WaitingForPayment), string(Paid):
		return true
	}
	return false
}

func PaymentMethodEnumValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(UponReceipt), string(Online):
		return true
	}
	return false
}

func OrderConditionsEnumValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(WithFitting), string(WithoutFitting):
		return true
	}
	return false
}

func ConvertFittingToBool(f OrderConditions) bool {
	return f == WithFitting
}

func PhoneValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	matched, err := regexp.MatchString(fall.VALID_PHONE, value)
	if err != nil {
		return false
	}
	return matched
}

type CreateOrderDto struct {
	PaymentMethod      PaymentMethodEnum `json:"payment_method" validate:"required,paymentMethodEnumValidation"`
	DeliveryPointId    int               `json:"delivery_point_id" validate:"required,min=1"`
	Conditions         OrderConditions   `json:"order_conditions" validate:"required,orderConditionsEnumValidation"`
	RecipientFirstname string            `json:"recipient_firstname" validate:"required,min=2"`
	RecipientLastname  string            `json:"recipient_lastname" validate:"required,min=2"`
	RecipientPhone     string            `json:"recipient_phone" validate:"required,phoneValidation"`
	ModelSizeIds       []int             `json:"model_size_ids" validate:"required,dive,min=1"`
}

type CreateOrderInput struct {
	DeliveryPrice      float64 `json:"delivery_price" validate:"required,min=0"`
	TotalPrice         float64
	ProductsPrice      float64
	TotalDiscount      float64
	RecipientFirstname string
	RecipientLastname  string
	Conditions         OrderConditions `json:"order_conditions" validate:"required,orderConditionsEnumValidation"`
	RecipientPhone     string
	PaymentMethod      PaymentMethodEnum `json:"payment_method" validate:"required,paymentMethodEnumValidation"`
	DeliveryPointId    int               `json:"delivery_point_id" validate:"required,min=1"`
	CartItems          []*CartItemModel
}
