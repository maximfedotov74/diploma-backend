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
		orderRouter.Post("/", func(c *fiber.Ctx) error { return nil })
	}
}
