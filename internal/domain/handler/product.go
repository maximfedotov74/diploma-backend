package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/generator"
)

type productService interface {
	CreateProduct(ctx context.Context, dto model.CreateProductDto) fall.Error             // +
	GetProductPage(ctx context.Context, slug string) (*model.ProductRelation, fall.Error) // +
	FindProductById(ctx context.Context, id int) (*model.Product, fall.Error)
	FindProductModelById(ctx context.Context, id int) (*model.ProductModel, fall.Error)
	CreateModel(ctx context.Context, dto model.CreateProductModelDto) fall.Error                //+
	AddPhoto(ctx context.Context, dto model.CreateProducModelImg) fall.Error                    // +
	RemovePhoto(ctx context.Context, photoId int) fall.Error                                    // +
	DeleteProduct(ctx context.Context, id int) fall.Error                                       // +
	DeleteProductModel(ctx context.Context, id int) fall.Error                                  // +
	UpdateProduct(ctx context.Context, dto model.UpdateProductDto, id int) fall.Error           // +
	UpdateProductModel(ctx context.Context, dto model.UpdateProductModelDto, id int) fall.Error // +
	FindModelsColored(ctx context.Context, id int) ([]model.ProductModelColors, fall.Error)     // +
	AdminGetProducts(ctx context.Context, page int, brandId *int, categoryId *int) (*model.AdminProductResponse, fall.Error)
	AdminGetProductModels(ctx context.Context, id int) ([]model.AdminProductModelRelation, fall.Error)
	GetCatalogModels(ctx context.Context, query generator.CatalogFilters) (*model.CatalogResponse, fall.Error)
	GetModelImages(ctx context.Context, modelId int) ([]model.ProductModelImg, fall.Error)
	GetModelSizes(ctx context.Context, modelId int) ([]model.ProductModelSize, fall.Error)
	GetModelOptions(ctx context.Context, modelId int) ([]*model.ProductModelOption, fall.Error)
}

type ProductHandler struct {
	service        productService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewProductHandler(service productService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *ProductHandler {
	return &ProductHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func (h *ProductHandler) InitRoutes() {
	productRouter := h.router.Group("product")
	{
		productRouter.Post("/", h.createProduct)
		productRouter.Post("/model", h.createProductModel)
		productRouter.Post("/model/img", h.addPhoto)

		productRouter.Delete("/model/img/:imgId", h.removePhoto)
		productRouter.Delete("/model/:modelId", h.deleteModel)
		productRouter.Delete("/:id", h.deleteProduct)

		productRouter.Patch("/model/:id", h.updateProductModel)
		productRouter.Patch("/:id", h.updateProduct)

		productRouter.Get("/admin", h.adminGetProducts)
		productRouter.Get("/admin/models/:productId", h.adminGetProductModels)
		productRouter.Get("/catalog/:categorySlug", h.getCatalogModels)

		productRouter.Get("/model/colors/:id", h.findModelsColored)
		productRouter.Get("/model/page/:slug", h.getProductPage)
		productRouter.Get("/model/img/:id", h.getProductModelImg)
		productRouter.Get("/model/sizes/:id", h.getProductModelSizes)
		productRouter.Get("/model/options/:id", h.getModelOptions)
	}
}

// @Summary Get catalog models
// @Description Get catalog models
// @Tags product
// @Accept json
// @Produce json
// @Param categorySlug path string true "Category slug"
// @Param size query string false "sizes"
// @Param brands query string false "brands"
// @Param sort query string false "sort by"
// @Param is_sale query string false "get items with sale"
// @Param price query string false "from - to"
// @Param page query string false "pagination page"
// @Router /api/product/catalog/{categorySlug} [get]
// @Success 200 {object} model.CatalogResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) getCatalogModels(ctx *fiber.Ctx) error {

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

	page, ok := query["page"]

	if ok {
		delete(query, "page")
	}

	filters := generator.CatalogFilters{
		Options:          query,
		Slug:             slug,
		Sizes:            sizes,
		Brands:           brands,
		SortBy:           sortBy,
		OnlyWithDiscount: onlyWithDiscount,
		Price:            price,
		Page:             page,
	}

	res, ex := h.service.GetCatalogModels(ctx.Context(), filters)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(res)
}

// @Summary Get product models
// @Security BearerToken
// @Description Get product models
// @Tags product
// @Accept json
// @Produce json
// @Param productId path int true "product id"
// @Router /api/product/admin/models/{productId} [get]
// @Success 200 {array} model.AdminProductModelRelation
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) adminGetProductModels(ctx *fiber.Ctx) error {

	productId, err := ctx.ParamsInt("productId")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	models, ex := h.service.AdminGetProductModels(ctx.Context(), productId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(models)
}

// @Summary Get product model options
// @Description Get product model options
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "model id"
// @Router /api/product/model/options/{id} [get]
// @Success 200 {array} model.ProductModelOption
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) getModelOptions(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	options, ex := h.service.GetModelOptions(ctx.Context(), id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(options)
}

// @Summary Get product model images
// @Description Get product model images
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "model id"
// @Router /api/product/model/img/{id} [get]
// @Success 200 {array} model.ProductModelImg
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) getProductModelImg(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	images, ex := h.service.GetModelImages(ctx.Context(), id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(images)
}

// @Summary Get product model sizes
// @Description Get product model sizes
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "model id"
// @Router /api/product/model/sizes/{id} [get]
// @Success 200 {array} model.ProductModelSize
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) getProductModelSizes(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	sizes, ex := h.service.GetModelSizes(ctx.Context(), id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(sizes)
}

// @Summary Get products for admin panel
// @Security BearerToken
// @Description Get products for admin panel
// @Tags product
// @Accept json
// @Produce json
// @Param page query int true "Page for pagination"
// @Param categoryId query int false "categoryId"
// @Param brandId query int false "brandId"
// @Router /api/product/admin [get]
// @Success 200 {object} model.AdminProductResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) adminGetProducts(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	var categoryId *int
	var brandId *int

	category := ctx.QueryInt("categoryId", 0)
	brand := ctx.QueryInt("brandId", 0)

	if brand > 0 {
		brandId = &brand
	}

	if category > 0 {
		categoryId = &category
	}

	prods, ex := h.service.AdminGetProducts(ctx.Context(), page, brandId, categoryId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(prods)
}

// @Summary Get product models color
// @Description Get product models color
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "product id"
// @Router /api/product/model/colors/{id} [get]
// @Success 200 {array} model.ProductModelColors
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) findModelsColored(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	models, ex := h.service.FindModelsColored(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(models)

}

// @Summary Update product model
// @Description Update product model
// @Tags product
// @Accept json
// @Produce json
// @Param dto body model.UpdateProductModelDto true "Update product model with body dto"
// @Param id path int true "product model id"
// @Router /api/product/model/{id} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) updateProductModel(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateProductModelDto{}

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

	ex := h.service.UpdateProductModel(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Update product
// @Description Update product
// @Tags product
// @Accept json
// @Produce json
// @Param dto body model.UpdateProductDto true "Update product with body dto"
// @Param id path int true "product id"
// @Router /api/product/{id} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) updateProduct(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	dto := model.UpdateProductDto{}

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

	ex := h.service.UpdateProduct(ctx.Context(), dto, id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete product
// @Description Delete product
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product id"
// @Router /api/product/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) deleteProduct(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteProduct(ctx.Context(), id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete product model
// @Description Delete product model
// @Tags product
// @Accept json
// @Produce json
// @Param modelId path int true "Product model id"
// @Router /api/product/model/{modelId} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) deleteModel(ctx *fiber.Ctx) error {

	modelId, err := ctx.ParamsInt("modelId")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.DeleteProductModel(ctx.Context(), modelId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete product model img
// @Description Delete product model img
// @Tags product
// @Accept json
// @Produce json
// @Param imgId path int true "Product model img id"
// @Router /api/product/model/img/{imgId} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) removePhoto(ctx *fiber.Ctx) error {

	imgId, err := ctx.ParamsInt("imgId")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := h.service.RemovePhoto(ctx.Context(), imgId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Get product model page
// @Description Get product model page
// @Tags product
// @Accept json
// @Produce json
// @Param slug path string true "Product Model Slug"
// @Router /api/product/model/page/{slug} [get]
// @Success 200 {object} model.ProductRelation
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) getProductPage(ctx *fiber.Ctx) error {

	slug := ctx.Params("slug")

	page, ex := h.service.GetProductPage(ctx.Context(), slug)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(page)
}

// @Summary Add img to product model
// @Description Add img to product model
// @Tags product
// @Accept json
// @Produce json
// @Param dto body model.CreateProducModelImg true "Add product model img with body dto"
// @Router /api/product/model/img [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) addPhoto(ctx *fiber.Ctx) error {
	dto := model.CreateProducModelImg{}

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

	ex := h.service.AddPhoto(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Create product
// @Description Create product
// @Tags product
// @Accept json
// @Produce json
// @Param dto body model.CreateProductDto true "Create product with body dto"
// @Router /api/product/ [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) createProduct(ctx *fiber.Ctx) error {
	dto := model.CreateProductDto{}

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

	ex := h.service.CreateProduct(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Create product model
// @Description Create product model
// @Tags product
// @Accept json
// @Produce json
// @Param dto body model.CreateProductModelDto true "Create product model with body dto"
// @Router /api/product/model/ [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *ProductHandler) createProductModel(ctx *fiber.Ctx) error {
	dto := model.CreateProductModelDto{}

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

	ex := h.service.CreateModel(ctx.Context(), dto)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}
