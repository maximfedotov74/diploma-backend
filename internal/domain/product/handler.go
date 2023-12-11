package product

import (
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Service interface {
	FindModelByIdWithRelations(id int) (*Product, exception.Error)
	CreateProduct(dto CreateProductDto) exception.Error
	CreateModel(dto CreateProductModelDto) exception.Error
	AddPhoto(dto CreateProducModelImg) exception.Error
	FindById(id int) (*ProductWithoutRelations, exception.Error)
	FindModelsColored(slug string) ([]ProductModelColors, exception.Error)
	AdminGetProducts(page int, brandId *int, categoryId *int) (*AdminProductResponse, exception.Error)
	RemovePhoto(photoId int) exception.Error
	GetCatalogModels(query utils.CatalogFilters) string
	UpdateProduct(dto UpdateProductDto, id int) exception.Error
	UpdateProductModel(dto UpdateProductModelDto, modelId int) exception.Error
	DeleteProduct(id int) exception.Error
	DeleteProductModel(id int) exception.Error
}

type RoleGuard func(roles ...string) fiber.Handler
type AuthGuard fiber.Handler

type ProductHandler struct {
	service   Service
	router    fiber.Router
	authGuard AuthGuard
	roleGuard RoleGuard
}

func NewProductHandler(service Service, router fiber.Router, authGuard AuthGuard, roleGuard RoleGuard) *ProductHandler {
	return &ProductHandler{service: service, router: router, authGuard: authGuard, roleGuard: roleGuard}
}

func (ph *ProductHandler) InitRoutes() {
	productRouter := ph.router.Group("product")
	{
		productRouter.Post("/", ph.createProduct)
		productRouter.Post("/model", ph.createProductModel)
		productRouter.Post("/add-photo", ph.addPhoto)
		productRouter.Get("/colored-models/:slug", ph.findModelsColored)
		productRouter.Get("/catalog-models/:categorySlug", ph.getCatalogModels)
		productRouter.Get("/admin/products-list", ph.adminGetProducts)
		productRouter.Get("/model-full/:id", ph.findModelByWithRelations)
		productRouter.Delete("/remove-photo/:photoId", ph.removePhoto)
		productRouter.Patch("/:id", ph.updateProduct)
		productRouter.Patch("/model/:id", ph.updateProductModel)
		productRouter.Delete("/:id", ph.deleteProduct)
		productRouter.Delete("/model/:id", ph.deleteProductModel)
	}
}

func (ph *ProductHandler) getCatalogModels(ctx *fiber.Ctx) error {

	slug := ctx.Params("categorySlug")

	query := ctx.Queries()

	sizes, ok := query["size"]

	if ok {
		delete(query, "size")
	}

	brands, ok := query["brands"]

	if ok {
		delete(query, "brands")
	}

	sortBy, ok := query["sort"]

	if ok {
		delete(query, "sort")
	}

	onlyWithDiscount, ok := query["is_sale"]

	if ok {
		delete(query, "is_sale")
	}

	price, ok := query["price"]

	if ok {
		delete(query, "price")
	}

	filters := utils.CatalogFilters{
		Options:          query,
		Slug:             slug,
		Sizes:            sizes,
		Brands:           brands,
		SortBy:           sortBy,
		OnlyWithDiscount: onlyWithDiscount,
		Price:            price,
	}

	generated := ph.service.GetCatalogModels(filters)

	return ctx.Status(exception.STATUS_OK).SendString(generated)
}

// @Summary Get color models
// @Description Get color models
// @Tags product
// @Accept json
// @Produce json
// @Param slug path string true "Product slug"
// @Router /api/colored-models/:slug [get]
// @Success 200 {array} product.ProductModelColors
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) findModelsColored(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	models, ex := ph.service.FindModelsColored(slug)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(models)
}

// @Summary Get products for admin panel
// @Security BearerToken
// @Description Get products for admin panel
// @Tags product
// @Accept json
// @Produce json
// @Param page query int true "Page for pagination"
// @Router /api/product/admin/products-list [get]
// @Success 201 {object} product.AdminProductResponse
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) adminGetProducts(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	var categoryId *int
	var brandId *int

	category := ctx.QueryInt("category", 0)
	brand := ctx.QueryInt("brand", 0)

	if brand > 0 {
		brandId = &brand
	}

	if category > 0 {
		categoryId = &category
	}

	prods, ex := ph.service.AdminGetProducts(page, brandId, categoryId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(prods)
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

	ex := ph.service.CreateProduct(dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
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

	ex := ph.service.CreateModel(dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Remove photo from product model
// @Description Remove photo from product model
// @Tags product
// @Accept json
// @Produce json
// @Param photoId path int true "Photo id"
// @Router /api/product/remove-photo/:photoId [delete]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) removePhoto(ctx *fiber.Ctx) error {
	photoId, err := ctx.ParamsInt("photoId")

	if err != nil {
		ex := exception.NewErr("PhotoId is required parameter or invalid format", exception.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	ex := ph.service.RemovePhoto(photoId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
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

	isValid := filepath.IsAbs(dto.ImgPath)

	if !isValid {
		ex := exception.NewErr("Некорректный путь до файла!", exception.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	ex := ph.service.AddPhoto(dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

// @Summary Get product model by model id
// @Description Get product model by model id
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "model id parameter"
// @Router /api/product/model-full/:id [get]
// @Success 200 {object} product.Product
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) findModelByWithRelations(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	p, ex := ph.service.FindModelByIdWithRelations(id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(exception.STATUS_OK).JSON(p)
}

// @Summary Update product
// @Description Update product
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "product id parameter"
// @Param dto body product.UpdateProductDto true "update product dto"
// @Router /api/product/:id [patch]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) updateProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := UpdateProductDto{}

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

	ex := ph.service.UpdateProduct(dto, id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Update product model
// @Description Update product model
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "model id parameter"
// @Param dto body product.UpdateProductModelDto true "update product-model dto"
// @Router /api/product/model/:id [patch]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) updateProductModel(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := UpdateProductModelDto{}

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

	ex := ph.service.UpdateProductModel(dto, id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Delete product
// @Description Delete product
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "product id parameter"
// @Router /api/product/:id [delete]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) deleteProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := ph.service.DeleteProduct(id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

// @Summary Delete product model
// @Description Delete product model
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "product model id parameter"
// @Router /api/product/model/:id [delete]
// @Success 200
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ph *ProductHandler) deleteProductModel(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := ph.service.DeleteProductModel(id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}
