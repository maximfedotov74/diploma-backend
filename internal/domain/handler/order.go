package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type orderService interface {
	Create(ctx context.Context, dto model.CreateOrderDto, user model.LocalSession) fall.Error
}

type OrderHandler struct {
	service        orderService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewOrderHandler(service orderService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *OrderHandler {
	return &OrderHandler{service: service, router: router, authMiddleware: authMiddleware}
}

func (h *OrderHandler) InitRoutes() {
	orderRouter := h.router.Group("order")
	{
		orderRouter.Post("/", h.create)
	}
}

// @Summary Create order
// @Description Create order
// @Tags order
// @Accept json
// @Produce json
// @Param dto body model.CreateOrderDto true "Create order with body dto"
// @Router /api/order/ [post]
// @Success 201 {object} model.Order
// @Failure 401 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OrderHandler) create(ctx *fiber.Ctx) error {
	return nil
}
