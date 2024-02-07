package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type categoryService interface {
	Create(ctx context.Context, dto model.CreateCategoryDto) fall.Error
	Update(ctx context.Context, dto model.UpdateCategoryDto, id int) fall.Error
	FindBySlug(ctx context.Context, slug string) (*model.CategoryModel, fall.Error)
	FindBySlugRelation(ctx context.Context, slug string) (*model.Category, fall.Error)
	Delete(ctx context.Context, slug string) fall.Error
	GetAll(ctx context.Context) ([]*model.Category, fall.Error)
	GetCatalogCategories(ctx context.Context, slug string) (*model.СatalogCategory, fall.Error)
	GetTopLevels(ctx context.Context) ([]model.CategoryModel, fall.Error)
}

type CategoryHandler struct {
	service        categoryService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewCategoryHandler(service categoryService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *CategoryHandler {
	return &CategoryHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func (h *CategoryHandler) InitRoutes() {
	categoryRouter := h.router.Group("category")
	{
		{
			categoryRouter.Get("/", h.getAll)
			categoryRouter.Get("/top", h.getTopLevels)
			categoryRouter.Get("/catalog/:slug", h.catalog)
			categoryRouter.Get("/relation/:slug", h.findBySlugRelation)
			categoryRouter.Post("/", h.create)
			categoryRouter.Patch("/:id", h.update)
			categoryRouter.Delete("/:slug", h.delete)
			categoryRouter.Get("/:slug", h.findBySlug)
		}
	}
}

// @Summary Get category by slug with subcategories
// @Description Get category by slug with subcategories
// @Tags category
// @Accept json
// @Produce json
// @Param slug path string true "Category Slug"
// @Router /api/category/relation/{slug} [get]
// @Success 200 {object} model.CategoryRelation
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) findBySlugRelation(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	category, err := h.service.FindBySlugRelation(ctx.Context(), slug)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(category)
}

// @Summary Get category by slug
// @Description Get category by slug
// @Tags category
// @Accept json
// @Produce json
// @Param slug path string true "Category Slug"
// @Router /api/category/{slug} [get]
// @Success 200 {object} model.CategoryModel
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) findBySlug(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	category, err := h.service.FindBySlug(ctx.Context(), slug)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(category)
}

// @Summary Get top level categories
// @Description Get top level categories
// @Tags category
// @Accept json
// @Produce json
// @Router /api/category/top [get]
// @Success 200 {array} model.CategoryModel
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) getTopLevels(ctx *fiber.Ctx) error {

	categories, err := h.service.GetTopLevels(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(categories)
}

// @Summary Get all categories
// @Description Get all categories
// @Tags category
// @Accept json
// @Produce json
// @Router /api/category/ [get]
// @Success 200 {array} model.CategoryRelation
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) getAll(ctx *fiber.Ctx) error {

	categories, err := h.service.GetAll(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(categories)
}

// @Summary Get catalog categories
// @Description Get catalog categories
// @Tags category
// @Accept json
// @Produce json
// @Param slug path string true "Category Slug"
// @Router /api/category/catalog/{slug} [get]
// @Success 200 {array} model.CatalogCategoryRelation
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) catalog(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	categories, err := h.service.GetCatalogCategories(ctx.Context(), slug)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(categories)
}

// @Summary Create category
// @Description Create category
// @Tags category
// @Accept json
// @Produce json
// @Param dto body model.CreateCategoryDto true "Create category with body dto"
// @Router /api/category/ [post]
// @Success 201
// @Failure 400 {object} fall.ValidationError
// @Failure 401 {object} fall.AppErr
// @Failure 403 {object} fall.AppErr
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) create(ctx *fiber.Ctx) error {
	dto := model.CreateCategoryDto{}

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

// @Summary Update category
// @Description Update category
// @Tags category
// @Accept json
// @Produce json
// @Param dto body model.UpdateCategoryDto true "Update category with body dto"
// @Param id path int true "Category Id"
// @Router /api/category/{id} [patch]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 401 {object} fall.AppErr
// @Failure 403 {object} fall.AppErr
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) update(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateCategoryDto{}

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

// @Summary Delete category by slug
// @Description Delete category by slug
// @Tags category
// @Accept json
// @Produce json
// @Param slug path string true "Category Slug"
// @Router /api/category/{slug} [delete]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *CategoryHandler) delete(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	err := h.service.Delete(ctx.Context(), slug)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.SendStatus(fall.STATUS_OK)
}
