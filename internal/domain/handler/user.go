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

type userService interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
	Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error
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
		userRouter.Get("/session", h.authMiddleware, h.getSession)
		userRouter.Get("/profile", h.authMiddleware, h.getProfile)
		userRouter.Patch("/profile", h.authMiddleware, h.updateProfile)
	}
}

// @Summary Get base profile info
// @Security BearerToken
// @Description Get base profile info
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/profile [get]
// @Success 200 {object} model.User
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) getProfile(ctx *fiber.Ctx) error {

	session, err := utils.GetLocalSession(ctx)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	user, err := h.service.FindById(ctx.Context(), session.UserId)
	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(user)
}

// @Summary Get local session
// @Security BearerToken
// @Description Get local session
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/session [get]
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

// @Summary Update user profile info
// @Description Update user profile info
// @Tags user
// @Accept json
// @Produce json
// @Param dto body model.UpdateUserDto true "Update user profile with body dto"
// @Router /api/user/profile [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) updateProfile(ctx *fiber.Ctx) error {

	session, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.UpdateUserDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	validate.RegisterValidation("userGenderEnumValidation", model.UserGenderEnumValidation)

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	ex = h.service.Update(ctx.Context(), dto, session.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)

}
