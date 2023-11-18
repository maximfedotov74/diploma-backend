package product

import (
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Service interface {
	FindByProductSlugAndModelId(slug string, id int) (*Product, exception.Error)
	CreateProduct(dto CreateProductDto) exception.Error
	CreateModel(dto CreateProductModelDto) exception.Error
	AddPhoto(dto CreateProducModelImg) exception.Error
	FindById(id int) (*ProductWithoutRelations, exception.Error)
}

type RoleGuard interface {
	CheckRoles(roles ...string) fiber.Handler
}

type ProductHandler struct {
	service   Service
	router    fiber.Router
	authGuard fiber.Handler
	roleGuard RoleGuard
}

func NewProductHandler(service Service, router fiber.Router, authGuard fiber.Handler, roleGuard RoleGuard) *ProductHandler {
	return &ProductHandler{service: service, router: router, authGuard: authGuard, roleGuard: roleGuard}
}

func (ph *ProductHandler) InitRoutes() {
	productRouter := ph.router.Group("product")
	{
		productRouter.Post("/", ph.createProduct)
		productRouter.Post("/model", ph.createProductModel)
		productRouter.Post("/add-photo", ph.addPhoto)
		productRouter.Get("/:id/:slug", ph.findByProductSlugAndModelId)
	}
}

// @Summary Create product
// @Description Create product with dto
// @Tags product
// @Accept json
// @Produce json
// @Param dto body product.CreateProductDto true "Create product dto"
// @Router /api/product/ [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) createProduct(ctx *fiber.Ctx) error {

	dto := CreateProductDto{}

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

	ex := ph.service.CreateProduct(dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(201)
}

// @Summary Create product-model
// @Description Create product model with dto
// @Tags product
// @Accept json
// @Produce json
// @Param dto body product.CreateProductModelDto true "Create product model dto"
// @Router /api/product/model [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) createProductModel(ctx *fiber.Ctx) error {

	dto := CreateProductModelDto{}

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

	ex := ph.service.CreateModel(dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(201)
}

// @Summary Add photo to product model
// @Description Add photo to product model
// @Tags product
// @Accept json
// @Produce json
// @Param dto body product.CreateProducModelImg true "Add photo to product model"
// @Router /api/product/add-photo [post]
// @Success 201
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) addPhoto(ctx *fiber.Ctx) error {

	dto := CreateProducModelImg{}

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

	isValid := filepath.IsAbs(dto.ImgPath)

	if !isValid {
		ex := exception.NewErr("Некорректный путь до файла!", 400)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	ex := ph.service.AddPhoto(dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(201)
}

// @Summary Get product by title
// @Description Get product by title
// @Tags product
// @Accept json
// @Produce json
// @Param slug path string true "product slug parameter"
// @Param id path int true "model id parameter"
// @Router /api/product/:slug [get]
// @Success 200 {object} product.Product
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) findByProductSlugAndModelId(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.INVALID_BODY, 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	p, ex := ph.service.FindByProductSlugAndModelId(slug, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(200).JSON(p)
}
