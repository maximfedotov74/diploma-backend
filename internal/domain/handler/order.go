package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type orderService interface {
	Create(ctx context.Context, dto model.CreateOrderDto, user *model.LocalSession) (*string, fall.Error)
	GetAdminOrders(ctx context.Context) ([]*model.Order, fall.Error)
	GetUserOrders(ctx context.Context, userId int) ([]*model.Order, fall.Error)
	GetOrder(ctx context.Context, id string) (*model.Order, fall.Error)
	CancelOrder(ctx context.Context, orderId string, userId int) fall.Error
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
		orderRouter.Patch("/cancel/:orderId", h.authMiddleware, h.cancel)
		orderRouter.Get("/my", h.authMiddleware, h.getUserOrders)
		orderRouter.Get("/:orderId", h.getOrder)
	}
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
