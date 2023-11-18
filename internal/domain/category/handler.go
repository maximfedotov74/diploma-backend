package category

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Service interface {
	CreateCategory(dto CreateCategoryDto) exception.Error
	GetCatalogCategories() ([]CatalogCategory, exception.Error)
	FindByIdWithSubcategories(id int) (*Category, exception.Error)
	FindBySlugWithSubcategories(slug string) (*Category, exception.Error)
	FindBySlug(slug string) (*CategoryDb, exception.Error)
	FindById(id int) (*CategoryDb, exception.Error)
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
		categoryRouter.Post("/create-category", ch.createCategory)
		categoryRouter.Get("/catalog", ch.getCatalogCategories)
		categoryRouter.Get("/recursive", ch.rg)
	}
}

func (h *CategoryHandler) rg(ctx *fiber.Ctx) error {
	err, j := h.service.FindBySlugWithSubcategories("men")
	if err != nil {
		return ctx.JSON(err)
	}
	return ctx.JSON(j)
}

// @Summary Create category
// @Description Create category with dto
// @Tags category
// @Accept json
// @Produce json
// @Param dto body category.CreateCategoryDto true "Create category dto"
// @Router /api/category/create-category [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *CategoryHandler) createCategory(ctx *fiber.Ctx) error {
	dto := CreateCategoryDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := exception.NewErr(messages.INVALID_BODY, 400)
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

	return ctx.SendStatus(201)
}

func (h *CategoryHandler) getCatalogCategories(ctx *fiber.Ctx) error {

	categories, appErr := h.service.GetCatalogCategories()

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(200).JSON(categories)

}
