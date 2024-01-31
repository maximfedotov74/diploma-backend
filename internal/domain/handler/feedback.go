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

type feedbackService interface {
	AddFeedback(ctx context.Context, dto model.AddFeedbackDto, userId int) fall.Error
	ToggleHidden(ctx context.Context, feedbackId int) fall.Error
	GetModelFeedback(ctx context.Context, modelId int, order string) (*model.ModelFeedbackResponse, fall.Error)
	GetAll(ctx context.Context, order string) ([]model.Feedback, fall.Error)
	DeleteFeedback(ctx context.Context, feedbackId int) fall.Error
}

type FeedbackHandler struct {
	service        feedbackService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewFeedbackHandler(service feedbackService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *FeedbackHandler {
	return &FeedbackHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func (fh *FeedbackHandler) InitRoutes() {
	feedbackRouter := fh.router.Group("feedback")
	{
		feedbackRouter.Post("/", fh.authMiddleware, fh.addFeedback)
		feedbackRouter.Delete("/:id", fh.deleteFeedback)
		feedbackRouter.Patch("/:id", fh.toggleHidden)
		feedbackRouter.Get("/model/:modelId", fh.getModelFeedback)
		feedbackRouter.Get("/", fh.getAll)
	}
}

// @Summary Toggle hidden feedback
// @Description Toggle hidden feedback
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path int true "feedback id"
// @Router /api/feedback/{id} [patch]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (fh *FeedbackHandler) toggleHidden(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := fh.service.ToggleHidden(ctx.Context(), id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(fall.STATUS_OK)
}

// @Summary Get all feedback
// @Description Get all feedback
// @Tags feedback
// @Accept json
// @Produce json
// @Router /api/feedback/ [get]
// @Param order query string false "Order [ASC | DESC]"
// @Success 200 {array} model.Feedback
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (fh *FeedbackHandler) getAll(ctx *fiber.Ctx) error {
	order := ctx.Query("order", "ASC")

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	feedback, ex := fh.service.GetAll(ctx.Context(), order)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(feedback)
}

// @Summary Delete feedback by id
// @Description Delete feedback by id
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path int true "Feedback Slug"
// @Router /api/feedback/{id} [delete]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (fh *FeedbackHandler) deleteFeedback(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := fh.service.DeleteFeedback(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(fall.STATUS_OK)
}

// @Summary Add feedback to model
// @Security BearerToken
// @Description Add feedback to model
// @Tags feedback
// @Accept json
// @Produce json
// @Param dto body model.AddFeedbackDto true "Add feedback with body dto"
// @Router /api/feedback/ [post]
// @Success 201
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (fh *FeedbackHandler) addFeedback(ctx *fiber.Ctx) error {

	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.AddFeedbackDto{}

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

		return ctx.Status(validError.Status).JSON(validError)
	}

	ex = fh.service.AddFeedback(ctx.Context(), dto, user.UserId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(fall.STATUS_CREATED)
}

// @Summary Get model feedback by modelId
// @Description Get model feedback by modelId
// @Tags feedback
// @Accept json
// @Produce json
// @Param modelId path int true "model Id"
// @Router /api/feedback/model/{modelId} [get]
// @Success 200 {object} model.ModelFeedbackResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *FeedbackHandler) getModelFeedback(ctx *fiber.Ctx) error {

	modelId, err := ctx.ParamsInt("modelId")
	order := ctx.Query("order", "ASC")

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	feedback, ex := h.service.GetModelFeedback(ctx.Context(), modelId, order)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(feedback)
}
