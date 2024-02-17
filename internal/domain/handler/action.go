package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type actionService interface {
	Create(ctx context.Context, dto model.CreateActionDto) fall.Error
	AddModel(ctx context.Context, dto model.AddModelToActionDto) fall.Error
	GetAll(ctx context.Context) ([]model.Action, fall.Error)
}

type ActionHandler struct {
	service        actionService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewActionHandler(service actionService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *ActionHandler {
	return &ActionHandler{service: service, router: router, authMiddleware: authMiddleware}
}

func (h *ActionHandler) InitRoutes() {
	actionRouter := h.router.Group("action")
	{
		actionRouter.Post("/", h.create)
		actionRouter.Post("/model", h.addModelToAction)
		actionRouter.Get("/", h.getAll)
	}
}

// @Summary Create action
// @Description Create action
// @Tags action
// @Accept json
// @Produce json
// @Param dto body model.CreateActionDto true "Create action with body dto"
// @Router /api/action/ [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) create(ctx *fiber.Ctx) error {
	dto := model.CreateActionDto{}

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

	ex := h.service.Create(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Add model to action
// @Description Add model to action
// @Tags action
// @Accept json
// @Produce json
// @Param dto body model.AddModelToActionDto true "Add model to action with body dto"
// @Router /api/action/model [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) addModelToAction(ctx *fiber.Ctx) error {
	dto := model.AddModelToActionDto{}

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

	ex := h.service.AddModel(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Get all actions
// @Description Get all actions
// @Tags action
// @Accept json
// @Produce json
// @Router /api/action/ [get]
// @Success 200 {array} model.Action
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) getAll(ctx *fiber.Ctx) error {

	actions, err := h.service.GetAll(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(actions)

}
