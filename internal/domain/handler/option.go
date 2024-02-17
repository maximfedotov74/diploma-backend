package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type optionService interface {
	GetCatalogFilters(ctx context.Context, categorySlug *string, brandSlug *string, actionId *string) (*model.CatalogFilters, fall.Error)
	GetAll(ctx context.Context) ([]*model.Option, fall.Error)
	FindOptionById(ctx context.Context, id int) (*model.Option, fall.Error)                       // +
	CreateOption(ctx context.Context, dto model.CreateOptionDto) fall.Error                       // +
	CreateSize(ctx context.Context, dto model.CreateSizeDto) fall.Error                           // +
	CreateValue(ctx context.Context, dto model.CreateOptionValueDto) fall.Error                   // +
	DeleteOption(ctx context.Context, id int) fall.Error                                          // +
	DeleteValue(ctx context.Context, id int) fall.Error                                           // +
	DeleteSize(ctx context.Context, id int) fall.Error                                            // +
	DeleteSizeFromProductModel(ctx context.Context, modelSizeId int) fall.Error                   // +
	DeleteOptionFromProductModel(ctx context.Context, productModelOptionId int) fall.Error        // +
	AddOptionToProductModel(ctx context.Context, dto model.AddOptionToProductModelDto) fall.Error // +
	AddSizeToProductModel(ctx context.Context, dto model.AddSizeToProductModelDto) fall.Error     // +
	UpdateOption(ctx context.Context, dto model.UpdateOptionDto, id int) fall.Error               // +
	UpdateOptionValue(ctx context.Context, dto model.UpdateOptionValueDto, id int) fall.Error     // +
	GetAllSizes(ctx context.Context) ([]model.Size, fall.Error)
}

type OptionHandler struct {
	service        optionService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewOptionHandler(service optionService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *OptionHandler {
	return &OptionHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func (h *OptionHandler) InitRoutes() {
	optionRouter := h.router.Group("characteristics")
	{
		optionRouter.Delete("/option/model/:id", h.deleteOptionFromProductModel)
		optionRouter.Delete("/size/model/:id", h.deleteSizeFromProductModel)
		optionRouter.Delete("/size/:id", h.deleteSize)
		optionRouter.Delete("/option/:id", h.deleteOption)
		optionRouter.Delete("/value/:id", h.deleteValue)
		optionRouter.Post("/option/model", h.addOptionToProductModel)
		optionRouter.Post("/size/model", h.addSizeToProductModel)
		optionRouter.Post("/option", h.createOption)
		optionRouter.Post("/value", h.createOptionValue)
		optionRouter.Post("/size", h.createSize)
		optionRouter.Patch("/option/:id", h.updateOption)
		optionRouter.Patch("/value/:id", h.updateOptionValue)
		optionRouter.Get("/catalog/:slug", h.getCatalogFilters)
		optionRouter.Get("/option", h.getAll)
		optionRouter.Get("/size", h.getAllSizes)

	}
}

// @Summary Get catalog filters
// @Description Get catalog filters
// @Tags characteristics
// @Accept json
// @Produce json
// @Param slug path string true "Category slug"
// @Router /api/characteristics/catalog/{slug} [get]
// @Success 200 {object} model.CatalogFilters
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) getCatalogFilters(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	filters, err := h.service.GetCatalogFilters(ctx.Context(), &slug, nil, nil)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(filters)
}

// @Summary Get all sizes
// @Description Get all sizes
// @Tags characteristics
// @Accept json
// @Produce json
// @Router /api/characteristics/size [get]
// @Success 200 {array} model.Size
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) getAllSizes(ctx *fiber.Ctx) error {

	sizes, err := h.service.GetAllSizes(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(sizes)

}

// @Summary Get all options
// @Description Get all options
// @Tags characteristics
// @Accept json
// @Produce json
// @Router /api/characteristics/option [get]
// @Success 200 {array} model.Option
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) getAll(ctx *fiber.Ctx) error {

	options, err := h.service.GetAll(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(options)

}

// @Summary Delete size from product model by id
// @Description Delete size from product model by id
// @Tags characteristics
// @Accept json
// @Produce json
// @Param id path int true "Product Model Option id"
// @Router /api/characteristics/size/model/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) deleteSizeFromProductModel(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteSizeFromProductModel(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete option from product model by id
// @Description Delete option from product model by id
// @Tags characteristics
// @Accept json
// @Produce json
// @Param id path int true "Product Model Option id"
// @Router /api/characteristics/option/model/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) deleteOptionFromProductModel(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteOptionFromProductModel(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Add size to product model
// @Description size option to product model
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.AddSizeToProductModelDto true "Add size to product model with body dto"
// @Router /api/characteristics/size/model [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) addSizeToProductModel(ctx *fiber.Ctx) error {

	dto := model.AddSizeToProductModelDto{}

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

	ex := h.service.AddSizeToProductModel(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Add option to product model
// @Description Add option to product model
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.AddOptionToProductModelDto true "Add option to product model with body dto"
// @Router /api/characteristics/option/model [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) addOptionToProductModel(ctx *fiber.Ctx) error {

	dto := model.AddOptionToProductModelDto{}

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

	ex := h.service.AddOptionToProductModel(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete size by id
// @Description Delete size by id
// @Tags characteristics
// @Accept json
// @Produce json
// @Param id path int true "Size id"
// @Router /api/characteristics/size/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) deleteSize(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteSize(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete option by id
// @Description Delete option by id
// @Tags characteristics
// @Accept json
// @Produce json
// @Param id path int true "Option id"
// @Router /api/characteristics/option/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) deleteOption(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteOption(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete option value by id
// @Description Delete option value by id
// @Tags characteristics
// @Accept json
// @Produce json
// @Param id path int true "Option id"
// @Router /api/characteristics/value/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) deleteValue(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteValue(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Create size
// @Description Create size
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.CreateSizeDto true "Create size with body dto"
// @Router /api/characteristics/size [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) createSize(ctx *fiber.Ctx) error {

	dto := model.CreateSizeDto{}

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

	ex := h.service.CreateSize(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Create option
// @Description Create option
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.CreateOptionDto true "Create option with body dto"
// @Router /api/characteristics/option [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) createOption(ctx *fiber.Ctx) error {

	dto := model.CreateOptionDto{}

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

	ex := h.service.CreateOption(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Create option value
// @Description Create option value
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.CreateOptionValueDto true "Create option value with body dto"
// @Router /api/characteristics/value [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) createOptionValue(ctx *fiber.Ctx) error {

	dto := model.CreateOptionValueDto{}

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

	ex := h.service.CreateValue(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Update option
// @Description Update option
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.UpdateOptionDto true "Update option with body dto"
// @Param id path int true "Option id"
// @Router /api/characteristics/option/{id} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) updateOption(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateOptionDto{}

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

	ex := h.service.UpdateOption(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Update option value
// @Description Update option value
// @Tags characteristics
// @Accept json
// @Produce json
// @Param dto body model.UpdateOptionValueDto true "Update option value with body dto"
// @Param id path int true "Value id"
// @Router /api/characteristics/value/{id} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *OptionHandler) updateOptionValue(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateOptionValueDto{}

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

	ex := h.service.UpdateOptionValue(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}
