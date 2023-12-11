package option

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Service interface {
	UpdateOption(dto UpdateOptionDto, id int) exception.Error
	CreateOption(dto CreateOptionDto) exception.Error
	GetById(id int) (*Option, exception.Error)
	CreateValue(dto CreateOptionValueDto) exception.Error
	DeleteOption(id int) exception.Error
	DeleteValue(id int) exception.Error
	AddOptionToProductModel(dto AddOptionToProductModelDto) exception.Error
	CreateSize(dto CreateSizeDto) exception.Error
	AddSizeToProductModel(dto AddSizeToProductModelDto) exception.Error
	GetCatalogFilters(categorySlug string) (*CatalogFilters, exception.Error)
	GetAll() ([]Option, exception.Error)
}

type RoleGuard func(roles ...string) fiber.Handler
type AuthGuard fiber.Handler

type OptionHandler struct {
	service   Service
	router    fiber.Router
	authGuard AuthGuard
}

func NewOptionHandler(service Service, router fiber.Router, authGuard AuthGuard) *OptionHandler {
	return &OptionHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
	}
}

func (oh *OptionHandler) InitRoutes() {
	optionRouter := oh.router.Group("option")
	{
		optionRouter.Post("/", oh.createOption)
		optionRouter.Get("/", oh.getAll)
		optionRouter.Post("/size", oh.createSize)
		optionRouter.Post("/value", oh.createValue)
		optionRouter.Post("/add-to-product-model", oh.addToProductModel)
		optionRouter.Post("/size/add-to-product-model", oh.addSizeToProductModel)
		optionRouter.Patch("/:id", oh.updateOption)
		optionRouter.Delete("/:id", oh.deleteOption)
		optionRouter.Delete("/value/:id", oh.deleteValue)
		optionRouter.Get("/catalog-filters/:categorySlug", oh.getCatalogFilters)
		optionRouter.Get("/:id", oh.getById)
	}
}

// @Summary Create size
// @Description Create size with dto
// @Tags option
// @Accept json
// @Produce json
// @Param dto body option.CreateSizeDto true "Create size dto"
// @Router /api/option/size [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) createSize(ctx *fiber.Ctx) error {
	dto := CreateSizeDto{}

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

	ex := oh.service.CreateSize(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Create option
// @Description Create option with dto
// @Tags option
// @Accept json
// @Produce json
// @Param dto body option.CreateOptionDto true "Create option dto"
// @Router /api/option/ [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) createOption(ctx *fiber.Ctx) error {
	dto := CreateOptionDto{}

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

	ex := oh.service.CreateOption(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Update option
// @Description Update option with dto
// @Tags option
// @Accept json
// @Produce json
// @Param dto body option.UpdateOptionDto true "Update option dto"
// @Router /api/option/:id [patch]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) updateOption(ctx *fiber.Ctx) error {
	dto := UpdateOptionDto{}
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
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

	ex := oh.service.UpdateOption(dto, id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Get option by id
// @Description Get option by id
// @Tags option
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/option/:id [get]
// @Success 200 {object} option.Option
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	opt, ex := oh.service.GetById(id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(opt)
}

// @Summary Get all options
// @Description Get all options
// @Tags option
// @Accept json
// @Produce json
// @Router /api/option/ [get]
// @Success 200 {array} option.Option
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) getAll(ctx *fiber.Ctx) error {

	opt, ex := oh.service.GetAll()

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(opt)
}

// @Summary Get get Catalog Filters by category slug
// @Description Get Catalog Filters
// @Tags option
// @Accept json
// @Produce json
// @Param categorySlug path string true "categorySlug parameter"
// @Router /api/option/catalog-filters/:categorySlug [get]
// @Success 200 {object} option.CatalogFilters
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) getCatalogFilters(ctx *fiber.Ctx) error {
	slug := ctx.Params("categorySlug")

	if slug == "" {
		appErr := exception.NewErr("CategorySlug is required query param", exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	filters, ex := oh.service.GetCatalogFilters(slug)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(filters)
}

// @Summary Create value
// @Description Create value with dto
// @Tags option
// @Accept json
// @Produce json
// @Param dto body option.CreateOptionValueDto true "Create option value dto"
// @Router /api/option/value [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) createValue(ctx *fiber.Ctx) error {
	dto := CreateOptionValueDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := exception.NewErr(messages.INVALID_BODY, exception.STATUS_OK)
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

	ex := oh.service.CreateValue(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Delete option by id
// @Description Delete option by id
// @Tags option
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/option/:id [delete]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) deleteOption(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := oh.service.DeleteOption(id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Delete option value by id
// @Description Delete option value by id
// @Tags option
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/option/value:id [delete]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) deleteValue(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := oh.service.DeleteValue(id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Add option to product model
// @Description Add option to product model
// @Tags option
// @Accept json
// @Produce json
// @Param dto body option.AddOptionToProductModelDto true "Add to product model dto"
// @Router /api/option/add-to-product-model [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) addToProductModel(ctx *fiber.Ctx) error {
	dto := AddOptionToProductModelDto{}

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

	ex := oh.service.AddOptionToProductModel(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Add size to product model
// @Description Add size to product model
// @Tags option
// @Accept json
// @Produce json
// @Param dto body option.AddSizeToProductModelDto true "Add size to product model dto"
// @Router /api/option/size/add-to-product-model [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (oh *OptionHandler) addSizeToProductModel(ctx *fiber.Ctx) error {
	dto := AddSizeToProductModelDto{}

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

	ex := oh.service.AddSizeToProductModel(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}
