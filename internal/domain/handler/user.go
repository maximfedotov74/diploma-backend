package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type userService interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
}

type UserHandler struct {
	service        userService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewUserHandler(service userService, router fiber.Router, authMiddleware middleware.AuthMiddleware,
) *UserHandler {
	return &UserHandler{service: service, router: router, authMiddleware: authMiddleware}
}

func (h *UserHandler) InitRoutes() {
	userRouter := h.router.Group("/user")
	{
		userRouter.Get("/", h.authMiddleware, h.getSession)
	}
}

// @Summary Get local session
// @Security BearerToken
// @Description Get local session
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/ [get]
// @Success 200 {object} model.LocalSession
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) getSession(ctx *fiber.Ctx) error {

	session, err := utils.GetLocalSession(ctx)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(session)
}
