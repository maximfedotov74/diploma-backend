package service

import (
	"context"
	"fmt"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/generator"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type productRepository interface {
	FindModelsColored(ctx context.Context, id int) ([]model.ProductModelColors, fall.Error)
	GetProductPage(ctx context.Context, slug string) (*model.ProductRelation, fall.Error)
	RemovePhoto(ctx context.Context, photoId int) fall.Error
	AddPhoto(ctx context.Context, dto model.CreateProducModelImg) fall.Error
	FindProductModelById(ctx context.Context, id int) (*model.ProductModel, fall.Error)
	FindProductModelBySlug(ctx context.Context, slug string) (*model.ProductModel, fall.Error)
	FindProductById(ctx context.Context, id int) (*model.Product, fall.Error)
	UpdateProductModel(ctx context.Context, dto model.UpdateProductModelDto, modelId int) fall.Error
	UpdateProduct(ctx context.Context, dto model.UpdateProductDto, id int) fall.Error
	DeleteProductModel(ctx context.Context, id int) fall.Error
	DeleteProduct(ctx context.Context, id int) fall.Error
	CreateModel(ctx context.Context, dto model.CreateProductModelDto, slug string) fall.Error
	CreateProduct(ctx context.Context, dto model.CreateProductDto) fall.Error
	AdminGetProducts(ctx context.Context, page int, brandId *int, categoryId *int) (*model.AdminProductResponse, fall.Error)
	AdminGetProductModels(ctx context.Context, id int) ([]model.ProductModel, fall.Error)
	GetCatalogModels(ctx context.Context, categorySlug string, sql generator.GeneratedCatalogQuery) (*model.CatalogResponse, fall.Error)
	GetModelImages(ctx context.Context, modelId int) ([]model.ProductModelImg, fall.Error)
	GetModelSizes(ctx context.Context, modelId int) ([]model.ProductModelSize, fall.Error)
	GetModelOptions(ctx context.Context, modelId int) ([]*model.ProductModelOption, fall.Error)
	SearchByArticle(ctx context.Context, article string) ([]model.SearchProductModel, fall.Error)
	GetSimilarProducts(ctx context.Context, filter model.SimilarProductsFilter, mainFilterId int,
		secondaryFilterId int, modelId int) ([]*model.CatalogProductModel, fall.Error)
	GetModelViews(ctx context.Context, modelId int) *int
	UpdateViews(ctx context.Context, ip string, modelId int)
	GetViewHistory(ctx context.Context, userId int, modelId int) ([]*model.CatalogProductModel, fall.Error)
	AddToViewHistory(ctx context.Context, userId int, modelId int) fall.Error
	GetPopularProducts(ctx context.Context, slug string) ([]*model.CatalogProductModel, fall.Error)
}

type productCategoryService interface {
	FindById(ctx context.Context, id int) (*model.CategoryModel, fall.Error)
	GetParentSubLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error)
	CheckForChildren(ctx context.Context, id int) (*int, fall.Error)
}

type productBrandService interface {
	FindById(ctx context.Context, id int) (*model.Brand, fall.Error)
}

type ProductService struct {
	repo            productRepository
	brandService    productBrandService
	categoryService productCategoryService
}

func NewProductService(repo productRepository, brandService productBrandService, categoryService productCategoryService) *ProductService {
	return &ProductService{
		repo:            repo,
		brandService:    brandService,
		categoryService: categoryService,
	}
}

func (s *ProductService) GetPopularProducts(ctx context.Context, slug string) ([]*model.CatalogProductModel, fall.Error) {
	return s.repo.GetPopularProducts(ctx, slug)
}

func (s *ProductService) AddToViewHistory(ctx context.Context, userId int, modelId int) fall.Error {
	return s.repo.AddToViewHistory(ctx, userId, modelId)
}

func (s *ProductService) GetViewHistory(ctx context.Context, userId int, modelId int) ([]*model.CatalogProductModel, fall.Error) {
	return s.repo.GetViewHistory(ctx, userId, modelId)
}

func (s *ProductService) GetSimilarProducts(ctx context.Context, categoryId int, brandId int, modelId int) ([]*model.CatalogProductModel, fall.Error) {
	bs, ex := s.repo.GetSimilarProducts(ctx, model.SimilarByBrand, brandId, categoryId, modelId)
	if ex != nil {
		return nil, ex
	}
	cs, ex := s.repo.GetSimilarProducts(ctx, model.SimilarByCategory, categoryId, brandId, modelId)
	if ex != nil {
		return nil, ex
	}

	res := append(bs, cs...)

	return res, nil
}

func (s *ProductService) FindProductModelBySlug(ctx context.Context, slug string) (*model.ProductModel, fall.Error) {
	return s.repo.FindProductModelBySlug(ctx, slug)
}

func (s *ProductService) SearchByArticle(ctx context.Context, article string) ([]model.SearchProductModel, fall.Error) {
	return s.repo.SearchByArticle(ctx, article)
}

func (s *ProductService) GetModelOptions(ctx context.Context, modelId int) ([]*model.ProductModelOption, fall.Error) {
	return s.repo.GetModelOptions(ctx, modelId)
}

func (s *ProductService) GetModelSizes(ctx context.Context, modelId int) ([]model.ProductModelSize, fall.Error) {
	return s.repo.GetModelSizes(ctx, modelId)
}

func (s *ProductService) GetModelImages(ctx context.Context, modelId int) ([]model.ProductModelImg, fall.Error) {
	return s.repo.GetModelImages(ctx, modelId)
}

func (s *ProductService) AdminGetProductModels(ctx context.Context, id int) ([]model.ProductModel, fall.Error) {
	return s.repo.AdminGetProductModels(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, dto model.CreateProductDto) fall.Error {
	cat, ex := s.categoryService.FindById(ctx, dto.CategoryId)
	if ex != nil {
		return ex
	}

	count, ex := s.categoryService.CheckForChildren(ctx, cat.Id)

	if ex != nil {
		return ex
	}

	if *count != 0 {
		return fall.NewErr(msg.ProductInvalidCategory, fall.STATUS_BAD_REQUEST)
	}

	_, ex = s.brandService.FindById(ctx, dto.BrandId)
	if ex != nil {
		return ex
	}
	return s.repo.CreateProduct(ctx, dto)
}

func (s *ProductService) CreateModel(ctx context.Context, dto model.CreateProductModelDto) fall.Error {
	p, ex := s.FindProductById(ctx, dto.ProductId)

	if ex != nil {
		return ex
	}

	categoryParent, ex := s.categoryService.GetParentSubLevel(ctx, p.Category.Id)
	if ex != nil {
		return ex
	}

	slug := utils.GenerateSlug(fmt.Sprintf("%s-%s-%s", categoryParent.Slug, p.Brand.Slug, p.Title))

	m, _ := s.repo.FindProductModelBySlug(ctx, slug)

	if m != nil {
		return fall.NewErr(msg.ProductModelSlugUnique, fall.STATUS_BAD_REQUEST)
	}

	return s.repo.CreateModel(ctx, dto, slug)
}

func (s *ProductService) GetProductPage(ctx context.Context, slug string) (*model.ProductRelation, fall.Error) {
	return s.repo.GetProductPage(ctx, slug)
}

func (s *ProductService) FindProductById(ctx context.Context, id int) (*model.Product, fall.Error) {
	return s.repo.FindProductById(ctx, id)
}

func (s *ProductService) FindProductModelById(ctx context.Context, id int) (*model.ProductModel, fall.Error) {
	return s.repo.FindProductModelById(ctx, id)
}

func (s *ProductService) AddPhoto(ctx context.Context, dto model.CreateProducModelImg) fall.Error {
	_, ex := s.FindProductModelById(ctx, dto.ProductModelId)
	if ex != nil {
		return ex
	}
	return s.repo.AddPhoto(ctx, dto)
}

func (s *ProductService) RemovePhoto(ctx context.Context, photoId int) fall.Error {
	return s.repo.RemovePhoto(ctx, photoId)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) fall.Error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) DeleteProductModel(ctx context.Context, id int) fall.Error {
	return s.repo.DeleteProductModel(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, dto model.UpdateProductDto, id int) fall.Error {
	p, ex := s.FindProductById(ctx, id)

	if ex != nil {
		return ex
	}

	return s.repo.UpdateProduct(ctx, dto, p.Id)
}

func (s *ProductService) UpdateProductModel(ctx context.Context, dto model.UpdateProductModelDto, id int) fall.Error {
	m, ex := s.FindProductModelById(ctx, id)
	if ex != nil {
		return ex
	}

	return s.repo.UpdateProductModel(ctx, dto, m.Id)
}

func (s *ProductService) FindModelsColored(ctx context.Context, id int) ([]model.ProductModelColors, fall.Error) {
	return s.repo.FindModelsColored(ctx, id)
}

func (s *ProductService) AdminGetProducts(ctx context.Context, page int, brandId *int, categoryId *int) (*model.AdminProductResponse, fall.Error) {
	return s.repo.AdminGetProducts(ctx, page, brandId, categoryId)
}

func (ps *ProductService) GetCatalogModels(ctx context.Context, query generator.CatalogFilters) (*model.CatalogResponse, fall.Error) {

	sql := generator.GenerateCatalogQuery(query)

	return ps.repo.GetCatalogModels(ctx, query.Slug, sql)
}

func (s *ProductService) UpdateViews(ctx context.Context, ip string, modelId int) {
	s.repo.UpdateViews(ctx, ip, modelId)
}

func (s *ProductService) GetModelViews(ctx context.Context, modelId int) *int {
	return s.repo.GetModelViews(ctx, modelId)
}
