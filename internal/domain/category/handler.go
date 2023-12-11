package category

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Service interface {
	CreateCategory(dto CreateCategoryDto) exception.Error
	FindByIdWithSubcategories(id int) (*Category, exception.Error)
	FindBySlugWithSubcategories(slug string) (*Category, exception.Error)
	FindBySlug(slug string) (*CategoryDb, exception.Error)
	FindById(id int) (*CategoryDb, exception.Error)
	GetAll() ([]Category, exception.Error)
	UpdateCategory(dto UpdateCategoryDto, id int) exception.Error
}

type CategoryHandler struct {
	service   Service
	router    fiber.Router
	authGuard fiber.Handler
}

func NewCategoryHandler(service Service, router fiber.Router, authGuard fiber.Handler) *CategoryHandler {
	return &CategoryHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
	}
}

func (ch *CategoryHandler) InitRoutes() {
	categoryRouter := ch.router.Group("/category")
	{
		categoryRouter.Get("/", ch.getAll)
		categoryRouter.Post("/", ch.createCategory)
		categoryRouter.Get("/with-sub/:slug", ch.getWithSub)
		categoryRouter.Patch("/:id", ch.updateCategory)
	}
}

// @Summary Update category
// @Description Update category
// @Tags category
// @Accept json
// @Produce json
// @Param dto body category.UpdateCategoryDto true "Update category dto"
// @Param id path int true "id parameter"
// @Router /api/category/:id [patch]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *CategoryHandler) updateCategory(ctx *fiber.Ctx) error {
	dto := UpdateCategoryDto{}

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(err.Error(), exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	err = ctx.BodyParser(&dto)

	if err != nil {
		appErr := exception.NewErr(messages.INVALID_BODY, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.service.UpdateCategory(dto, id)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Get all categories
// @Description Get all categories
// @Tags category
// @Accept json
// @Produce json
// @Router /api/category/ [get]
// @Success 200 {array} category.Category
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *CategoryHandler) getAll(ctx *fiber.Ctx) error {
	cats, ex := h.service.GetAll()
	if ex != nil {
		return ctx.JSON(ex)
	}
	return ctx.JSON(cats)
}

// @Summary Get by slug with sub categories
// @Description Get by slug with sub categories
// @Tags category
// @Accept json
// @Produce json
// @Router /api/category/with-sub/:slug [get]
// @Param slug path string true "slug parameter"
// @Success 200 {object} category.Category
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *CategoryHandler) getWithSub(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	cat, ex := h.service.FindBySlugWithSubcategories(slug)
	if ex != nil {
		return ctx.JSON(ex)
	}
	return ctx.JSON(cat)
}

// @Summary Create category
// @Description Create category with dto
// @Tags category
// @Accept json
// @Produce json
// @Param dto body category.CreateCategoryDto true "Create category dto"
// @Router /api/category/ [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *CategoryHandler) createCategory(ctx *fiber.Ctx) error {
	dto := CreateCategoryDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := exception.NewErr(messages.INVALID_BODY, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.service.CreateCategory(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}
