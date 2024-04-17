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
	Update(ctx context.Context, dto model.UpdateActionDto, id string) fall.Error
	GetModels(ctx context.Context, id string) ([]model.ActionModel, fall.Error)
	DeleteActionModel(ctx context.Context, actionModelId int) fall.Error
	DeleteAction(ctx context.Context, id string) fall.Error
	GetActionsByGender(ctx context.Context, gender model.ActionGender) ([]model.Action, fall.Error)
	FindById(ctx context.Context, id string) (*model.Action, fall.Error)
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
		actionRouter.Get("/by-id/:id", h.getById)
		actionRouter.Get("/by-gender/:gender", h.getByGender)
		actionRouter.Get("/model/:id", h.getActionModels)
		actionRouter.Patch("/:id", h.update)
		actionRouter.Delete("/model/:actionModelId", h.deleteActionModel)
		actionRouter.Delete("/:id", h.deleteAction)
	}
}

// @Summary Get action by id
// @Description Get action by id
// @Tags action
// @Accept json
// @Produce json
// @Param id path string true "action id"
// @Router /api/action/by-id/{id} [get]
// @Success 200 {object} model.Action
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) getById(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	if id == "" {
		ex := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)

	}

	action, ex := h.service.FindById(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(action)

}

// @Summary Get actions by gender
// @Description Get actions by gender
// @Tags action
// @Accept json
// @Produce json
// @Param gender path string true "action gender"
// @Router /api/action/by-gender/{gender} [get]
// @Success 200 {array} model.Action
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) getByGender(ctx *fiber.Ctx) error {

	gender := ctx.Params("gender")

	validate := validator.New()

	validate.RegisterValidation("actionGenderEnumValidation", model.ActionGenderEnumValidation)

	err := validate.Var(gender, "actionGenderEnumValidation")

	if err != nil {
		ex := fall.NewErr(err.Error(), fall.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	actions, ex := h.service.GetActionsByGender(ctx.Context(), model.ActionGender(gender))

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(actions)

}

// @Summary Delete action
// @Description Delete action
// @Tags action
// @Accept json
// @Produce json
// @Param id path string true "Action id"
// @Router /api/action/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) deleteAction(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	if id == "" {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteAction(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete action model
// @Description Delete action model
// @Tags action
// @Accept json
// @Produce json
// @Param actionModelId path int true "Action model  id"
// @Router /api/action/model/{actionModelId} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) deleteActionModel(ctx *fiber.Ctx) error {

	actionModelId, err := ctx.ParamsInt("actionModelId")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteActionModel(ctx.Context(), actionModelId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Update action
// @Description Update action
// @Tags action
// @Accept json
// @Produce json
// @Param dto body model.UpdateActionDto true "Update action with body dto"
// @Param id path string true "Action id"
// @Router /api/action/{id} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) update(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	if id == "" {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateActionDto{}

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

	ex := h.service.Update(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)

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

	validate.RegisterValidation("actionGenderEnumValidation", model.ActionGenderEnumValidation)

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

// @Summary Get action models
// @Description Get action models
// @Tags action
// @Accept json
// @Produce json
// @Router /api/action/model/{id} [get]
// @Param id path string true "Action id"
// @Success 200 {array} model.ActionModel
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ActionHandler) getActionModels(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	if id == "" {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	models, err := h.service.GetModels(ctx.Context(), id)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(models)

}
