package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type brandService interface {
	FindBySlug(ctx context.Context, slug string) (*model.Brand, fall.Error)
	GetAll(ctx context.Context) ([]model.Brand, fall.Error)
	Update(ctx context.Context, dto model.UpdateBrandDto, id int) fall.Error
	Create(ctx context.Context, dto model.CreateBrandDto) fall.Error
	Delete(ctx context.Context, slug string) fall.Error
}

type BrandHandler struct {
	service        brandService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewBrandHandler(service brandService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *BrandHandler {
	return &BrandHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func (h *BrandHandler) InitRoutes() {
	brandRouter := h.router.Group("brand")
	{
		brandRouter.Post("/", h.create)
		brandRouter.Patch("/:id", h.update)
		brandRouter.Delete("/:slug", h.delete)
		brandRouter.Get("/", h.getAll)
		brandRouter.Get("/:slug", h.findBySlug)
	}
}

// @Summary Delete brand by slug
// @Description Delete brand by slug
// @Tags brand
// @Accept json
// @Produce json
// @Param slug path string true "Brand Slug"
// @Router /api/brand/{slug} [delete]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *BrandHandler) delete(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	err := h.service.Delete(ctx.Context(), slug)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.SendStatus(fall.STATUS_OK)
}

// @Summary Get brand by slug
// @Description Get brand by slug
// @Tags brand
// @Accept json
// @Produce json
// @Param slug path string true "Brand Slug"
// @Router /api/brand/{slug} [get]
// @Success 200 {object} model.Brand
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *BrandHandler) findBySlug(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	brand, err := h.service.FindBySlug(ctx.Context(), slug)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(brand)
}

// @Summary Get all brands
// @Description Get all brands
// @Tags brand
// @Accept json
// @Produce json
// @Router /api/brand/ [get]
// @Success 200 {array} model.Brand
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *BrandHandler) getAll(ctx *fiber.Ctx) error {

	brands, err := h.service.GetAll(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(brands)

}

// @Summary Create brand
// @Description Create brand
// @Tags brand
// @Accept json
// @Produce json
// @Param dto body model.CreateBrandDto true "Create brand with body dto"
// @Router /api/brand/ [post]
// @Success 201
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *BrandHandler) create(ctx *fiber.Ctx) error {
	dto := model.CreateBrandDto{}

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

	return ctx.SendStatus(fall.STATUS_CREATED)

}

// @Summary Update brand
// @Description Update brand
// @Tags brand
// @Accept json
// @Produce json
// @Param dto body model.UpdateBrandDto true "Update brand with body dto"
// @Param id path int true "Brand id"
// @Router /api/brand/{id} [patch]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *BrandHandler) update(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateBrandDto{}

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

	ex := h.service.Update(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(fall.STATUS_OK)

}
