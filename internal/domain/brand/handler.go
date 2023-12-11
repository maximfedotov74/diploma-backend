package brand

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Service interface {
	CreateBrand(dto CreateBrandDto) exception.Error
	FindBySlug(slug string) (*Brand, exception.Error)
	GetAll() ([]Brand, exception.Error)
	UpdateBrand(dto UpdateBrandDto, id int) exception.Error
}

type BrandHandler struct {
	service Service
	router  fiber.Router
}

func NewBrandHandler(service Service, router fiber.Router) *BrandHandler {
	return &BrandHandler{service: service, router: router}
}

func (bh *BrandHandler) InitRoutes() {
	brandRouter := bh.router.Group("brand")
	{
		brandRouter.Post("/", bh.createBrand)
		brandRouter.Patch("/:id", bh.updateBrand)
		brandRouter.Get("/", bh.getAll)
		brandRouter.Get("/:slug", bh.findBySlug)
	}
}

// @Summary Update brand
// @Description Update brand
// @Tags brand
// @Accept json
// @Produce json
// @Param dto body brand.UpdateBrandDto true "Update brand dto"
// @Param id path int true "id parameter"
// @Router /api/brand/:id [patch]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *BrandHandler) updateBrand(ctx *fiber.Ctx) error {
	dto := UpdateBrandDto{}

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

	appErr := h.service.UpdateBrand(dto, id)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Create brand
// @Description Create brand with dto
// @Tags brand
// @Accept json
// @Produce json
// @Param dto body brand.CreateBrandDto true "Create brand dto"
// @Router /api/brand/ [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (bh *BrandHandler) createBrand(ctx *fiber.Ctx) error {
	dto := CreateBrandDto{}

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

	ex := bh.service.CreateBrand(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Get brand by id
// @Description Get brand by id
// @Tags brand
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/brand/:id [get]
// @Success 200 {object} brand.Brand
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (bh *BrandHandler) findBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	brand, ex := bh.service.FindBySlug(slug)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(exception.STATUS_OK).JSON(brand)
}

// @Summary Get all brands
// @Description Get all brands
// @Tags brand
// @Accept json
// @Produce json
// @Router /api/brand/ [get]
// @Success 200 {array} brand.Brand
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (bh *BrandHandler) getAll(ctx *fiber.Ctx) error {
	brands, ex := bh.service.GetAll()
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(exception.STATUS_OK).JSON(brands)
}
