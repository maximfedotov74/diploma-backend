package handler

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type orderService interface {
	Create(ctx context.Context, dto model.CreateOrderDto, user *model.LocalSession) (*string, fall.Error)
	GetAdminOrders(ctx context.Context, page int, fromDate *string, toDate *string) (*model.AllOrdersResponse, fall.Error)
	GetUserOrders(ctx context.Context, userId int) ([]*model.Order, fall.Error)
	GetOrder(ctx context.Context, id string) (*model.Order, fall.Error)
	CancelOrder(ctx context.Context, orderId string, userId int) fall.Error
	ChangeStatus(ctx context.Context, orderId string, status model.OrderStatusEnum) fall.Error
	ChangeDeliveryDate(ctx context.Context, orderId string, date time.Time) fall.Error
}

type OrderHandler struct {
	service        orderService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
	clientUrl      string
}

func NewOrderHandler(service orderService, router fiber.Router, authMiddleware middleware.AuthMiddleware, clientUrl string,
) *OrderHandler {
	return &OrderHandler{service: service, router: router, authMiddleware: authMiddleware, clientUrl: clientUrl}
}

func (h *OrderHandler) InitRoutes() {
	orderRouter := h.router.Group("order")
	{
		orderRouter.Post("/", h.authMiddleware, h.create)
		orderRouter.Get("/confirm-online-payment/:orderId", h.confirm)
		orderRouter.Get("/all", h.getAllOrders)
		orderRouter.Patch("/cancel/:orderId", h.authMiddleware, h.cancel)
		orderRouter.Get("/my", h.authMiddleware, h.getUserOrders)
		orderRouter.Get("/:orderId", h.getOrder)
		orderRouter.Patch("/change-status/:orderId", h.changeStatus)
		orderRouter.Patch("/change-delivery-date/:orderId", h.changeDeliveryDate)
	}
}

// @Summary Change order delivery date
// @Description Change order delivery date
// @Tags order
// @Accept json
// @Produce json
// @Param orderId path string true "Order id"
// @Param dto body model.ChangeOrderDeliveryDate true "Change delivery date with body dto"
// @Router /api/order/change-delivery-date/{orderId} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 401 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) changeDeliveryDate(ctx *fiber.Ctx) error {
	orderId := ctx.Params("orderId")

	dto := model.ChangeOrderDeliveryDate{}

	err := ctx.BodyParser(&dto)

	if err != nil {

		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	ex := h.service.ChangeDeliveryDate(ctx.Context(), orderId, dto.Date)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	ok := fall.GetOk()
	return ctx.Status(ok.Status()).JSON(ok)
}

// @Summary Change order status
// @Description Change order status
// @Tags order
// @Accept json
// @Produce json
// @Param orderId path string true "Order id"
// @Param dto body model.ChangeOrderStatusDto true "Change order status with body dto"
// @Router /api/order/change-status/{orderId} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 401 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) changeStatus(ctx *fiber.Ctx) error {
	orderId := ctx.Params("orderId")

	dto := model.ChangeOrderStatusDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {

		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	ex := h.service.ChangeStatus(ctx.Context(), orderId, dto.Status)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	ok := fall.GetOk()
	return ctx.Status(ok.Status()).JSON(ok)
}

// @Summary Get all orders
// @Security BearerToken
// @Description Get all orders
// @Tags order
// @Accept json
// @Produce json
// @Router /api/order/all [get]
// @Param fromDate query string false "from date"
// @Param toDate query string false "to date"
// @Param page query int false "pagination page"
// @Success 200 {object} model.AllOrdersResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) getAllOrders(ctx *fiber.Ctx) error {

	page := ctx.QueryInt("page", 1)

	var fromDate *string
	var toDate *string

	ISOlayout := "2006-01-02T15:04:05Z07:00"
	layout := "2006-01-02 15:04:05"

	fromQ := ctx.Query("fromDate")
	toQ := ctx.Query("toDate")

	if fromQ != "" {
		parsed, err := time.Parse(ISOlayout, fromQ)
		if err == nil {
			formatted := parsed.Format(layout)
			fromDate = &formatted
		}

	}
	if toQ != "" {
		parsed, err := time.Parse(ISOlayout, toQ)
		if err == nil {
			formatted := parsed.Format(layout)

			toDate = &formatted
		}
	}

	orders, ex := h.service.GetAdminOrders(ctx.Context(), page, fromDate, toDate)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(fall.STATUS_OK).JSON(orders)
}

// @Summary Cancel order
// @Description Cancel order
// @Tags order
// @Accept json
// @Produce json
// @Param orderId path string true "Order id"
// @Router /api/order/cancel/{orderId} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 401 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) cancel(ctx *fiber.Ctx) error {
	orderId := ctx.Params("orderId")

	user, ex := utils.GetLocalSession(ctx)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	ex = h.service.CancelOrder(ctx.Context(), orderId, user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	ok := fall.GetOk()
	return ctx.Status(ok.Status()).JSON(ok)
}

func (h *OrderHandler) confirm(ctx *fiber.Ctx) error {
	return ctx.Redirect(h.clientUrl, fall.STATUS_REDIRECT_PERM)
}

// @Summary Get order by id
// @Description Get order by id
// @Tags order
// @Accept json
// @Produce json
// @Router /api/order/{orderId} [get]
// @Param orderId path string true "Order id"
// @Success 200 {object} model.Order
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) getOrder(ctx *fiber.Ctx) error {

	orderId := ctx.Params("orderId")

	order, ex := h.service.GetOrder(context.Background(), orderId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(fall.STATUS_OK).JSON(order)
}

// @Summary Get user orders
// @Security BearerToken
// @Description Get user orders
// @Tags order
// @Accept json
// @Produce json
// @Router /api/order/my [get]
// @Success 200 {array} model.Order
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) getUserOrders(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	orders, ex := h.service.GetUserOrders(context.Background(), user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(fall.STATUS_OK).JSON(orders)
}

// @Summary Create order
// @Description Create order
// @Tags order
// @Accept json
// @Produce json
// @Param dto body model.CreateOrderDto true "Create order with body dto"
// @Router /api/order/ [post]
// @Success 201 {object} model.OrderConfirmation
// @Failure 401 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) create(ctx *fiber.Ctx) error {
	claims, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.CreateOrderDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	validate.RegisterValidation("paymentMethodEnumValidation", model.PaymentMethodEnumValidation)
	validate.RegisterValidation("orderConditionsEnumValidation", model.OrderConditionsEnumValidation)

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}
	link, ex := h.service.Create(ctx.Context(), dto, claims)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	r := model.OrderConfirmation{}

	if link != nil {
		r.PaymentUrl = link
		return ctx.Status(fall.STATUS_CREATED).JSON(r)
	}

	return ctx.Status(fall.STATUS_CREATED).JSON(r)
}
