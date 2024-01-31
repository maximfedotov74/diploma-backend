package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type deliveryRepository interface {
	Create(ctx context.Context, dto model.CreateDeliveryPointDto) fall.Error
	SearchPoints(ctx context.Context, text string, withFitting bool) ([]model.DeliveryPoint, fall.Error)
	FindById(ctx context.Context, id int) (*model.DeliveryPoint, fall.Error)
	Update(ctx context.Context, dto model.UpdateDeliveryPointDto, id int) fall.Error
	Delete(ctx context.Context, id int) fall.Error
}

type DeliveryHandler struct {
	repo           deliveryRepository
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewDeliveryHandler(repo deliveryRepository, r fiber.Router, m middleware.AuthMiddleware) *DeliveryHandler {
	return &DeliveryHandler{repo: repo, router: r, authMiddleware: m}
}

func (h *DeliveryHandler) InitRoutes() {
	deliveryRouter := h.router.Group("delivery")
	{
		deliveryRouter.Post("/", h.сreate)
		deliveryRouter.Patch("/:id", h.update)
		deliveryRouter.Delete("/:id", h.delete)
		deliveryRouter.Get("/search", h.search)
		deliveryRouter.Get("/:id", h.findById)
	}
}

// @Summary Create delivery-point
// @Description Create delivery-point
// @Tags delivery
// @Accept json
// @Produce json
// @Param dto body model.CreateDeliveryPointDto true "Create delivery-point with body dto"
// @Router /api/delivery/ [post]
// @Success 201
// @Failure 401 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *DeliveryHandler) сreate(ctx *fiber.Ctx) error {
	dto := model.CreateDeliveryPointDto{}

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

	ex := h.repo.Create(ctx.Context(), dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(fall.STATUS_CREATED)
}

// @Summary Update delivery-point
// @Description Update delivery-point
// @Tags delivery
// @Accept json
// @Produce json
// @Param dto body model.UpdateDeliveryPointDto true "Update delivery-point with body dto"
// @Param id path int true "delivery-point id"
// @Router /api/delivery/{id} [patch]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *DeliveryHandler) update(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateDeliveryPointDto{}

	err = ctx.BodyParser(&dto)

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

	ex := h.repo.Update(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(fall.STATUS_OK)

}

// @Summary Find delivery-point by id
// @Description Find delivery-point by id
// @Tags delivery
// @Accept json
// @Produce json
// @Param id path int true "delivery-point id"
// @Router /api/delivery/{id} [get]
// @Success 200 {object} model.DeliveryPoint
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *DeliveryHandler) findById(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	p, ex := h.repo.FindById(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(p)
}

// @Summary Delete delivery-point by id
// @Description Delete delivery-point by id
// @Tags delivery
// @Accept json
// @Produce json
// @Param id path int true "delivery-point id"
// @Router /api/delivery/{id} [delete]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *DeliveryHandler) delete(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.repo.Delete(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(fall.STATUS_OK)
}

// @Summary Find delivery-point by id
// @Description Find delivery-point by id
// @Tags delivery
// @Accept json
// @Produce json
// @Param with_fitting query bool false "delivery-point with_fitting filter"
// @Param search_text query string false "delivery-point search_text filter"
// @Router /api/delivery/search [get]
// @Success 200 {array} model.DeliveryPoint
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *DeliveryHandler) search(ctx *fiber.Ctx) error {

	with_fiitng := ctx.QueryBool("with_fitting")
	search_text := ctx.Query("search_text")

	p, ex := h.repo.SearchPoints(ctx.Context(), search_text, with_fiitng)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(p)
}
