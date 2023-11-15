package brand

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Service interface {
	CreateBrand(dto CreateBrandDto) exception.Error
	FindByTitle(title string) (*Brand, exception.Error)
	FindById(id int) (*Brand, exception.Error)
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
		brandRouter.Get("/:title", bh.findByTitle)
	}
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

	ex := bh.service.CreateBrand(dto)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(201)
}

// @Summary Get brand by title
// @Description Get brand by title
// @Tags brand
// @Accept json
// @Produce json
// @Param title path string true "title parameter"
// @Router /api/brand/:title [get]
// @Success 200 {object} brand.Brand
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (bh *BrandHandler) findByTitle(ctx *fiber.Ctx) error {
	title := ctx.Params("title")
	brand, ex := bh.service.FindByTitle(title)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(200).JSON(brand)
}
