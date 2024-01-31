package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type userService interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
}

type UserHandler struct {
	service userService
	router  fiber.Router
}

func NewUserHandler(service userService, router fiber.Router) *UserHandler {
	return &UserHandler{service: service, router: router}
}

func (h *UserHandler) InitRoutes() {
	userRouter := h.router.Group("/user")
	{
		userRouter.Get("/", func(c *fiber.Ctx) error {
			return nil
		})
	}
}
