package product

import (
	"fmt"

	"github.com/maximfedotov74/fiber-psql/internal/domain/brand"
	"github.com/maximfedotov74/fiber-psql/internal/domain/category"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Repository interface {
	CreateProduct(dto CreateProductDto, slug string) exception.Error
	FindModelByIdWithRelations(modelId int) (*Product, exception.Error)
	CreateModel(dto CreateProductModelDto) exception.Error
	AddPhoto(dto CreateProducModelImg) exception.Error
	FindById(id int) (*ProductWithoutRelations, exception.Error)
	FindModelById(modelId int) (*ProductModelWithoutRelations, exception.Error)
	FindModelsColored(slug string) ([]ProductModelColors, exception.Error)
	AdminGetProducts(page int, brandId *int, categoryId *int) (*AdminProductResponse, exception.Error)
	UpdateProduct(dto UpdateProductDto, slug *string, id int) exception.Error
	UpdateProductModel(dto UpdateProductModelDto, modelId int) exception.Error
	RemovePhoto(photoId int) exception.Error
	DeleteProductModel(id int) exception.Error
	DeleteProduct(id int) exception.Error
	GetCatalogModels(categorySlug string, sql utils.GeneratedCatalogQuery) string
}

type CategoryService interface {
	FindById(id int) (*category.CategoryDb, exception.Error)
	GetParentSubLevel(id int) (*category.CategoryDb, exception.Error)
}

type BrandService interface {
	FindById(id int) (*brand.Brand, exception.Error)
}

type ProductService struct {
	repo            Repository
	categoryService CategoryService
	brandService    BrandService
}

func NewProductService(repo Repository, cs CategoryService, bs BrandService) *ProductService {
	return &ProductService{
		repo:            repo,
		categoryService: cs,
		brandService:    bs,
	}
}

func (ps *ProductService) RemovePhoto(photoId int) exception.Error {
	err := ps.repo.RemovePhoto(photoId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) FindModelsColored(slug string) ([]ProductModelColors, exception.Error) {
	models, err := ps.repo.FindModelsColored(slug)
	if err != nil {
		return nil, err
	}

	return models, nil
}
func (ps *ProductService) GetCatalogModels(query utils.CatalogFilters) string {

	sql := utils.GenerateCatalogQuery(query)

	str := ps.repo.GetCatalogModels(query.Slug, sql)

	return str
}

func (ps *ProductService) AdminGetProducts(page int, brandId *int, categoryId *int) (*AdminProductResponse, exception.Error) {
	prods, err := ps.repo.AdminGetProducts(page, brandId, categoryId)
	if err != nil {
		return nil, err
	}

	return prods, nil
}

func (ps *ProductService) CreateProduct(dto CreateProductDto) exception.Error {

	category, ex := ps.categoryService.FindById(dto.CategoryID)
	if ex != nil {
		return ex
	}

	categoryParent, ex := ps.categoryService.GetParentSubLevel(category.Id)

	if ex != nil {
		return ex
	}

	brand, ex := ps.brandService.FindById(dto.BrandID)
	if ex != nil {
		return ex
	}

	slug := utils.GenerateSlug(fmt.Sprintf("%s-%s-%s", categoryParent.Slug, brand.Slug, dto.Title))

	ex = ps.repo.CreateProduct(dto, slug)
	if ex != nil {
		return ex
	}

	return nil
}

func (ps *ProductService) FindModelByIdWithRelations(id int) (*Product, exception.Error) {
	p, err := ps.repo.FindModelByIdWithRelations(id)

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (ps *ProductService) FindById(id int) (*ProductWithoutRelations, exception.Error) {
	p, err := ps.repo.FindById(id)

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (ps *ProductService) CreateModel(dto CreateProductModelDto) exception.Error {
	_, ex := ps.FindById(dto.ProductId)

	if ex != nil {
		return ex
	}

	ex = ps.repo.CreateModel(dto)

	if ex != nil {
		return ex
	}

	return nil
}

func (ps *ProductService) AddPhoto(dto CreateProducModelImg) exception.Error {
	err := ps.repo.AddPhoto(dto)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) UpdateProduct(dto UpdateProductDto, id int) exception.Error {
	product, err := ps.FindById(id)

	if err != nil {
		return err
	}

	categoryParent, ex := ps.categoryService.GetParentSubLevel(product.Category.Id)

	if ex != nil {
		return ex
	}

	var slug *string

	if dto.Title != nil {
		newSlug := utils.GenerateSlug(fmt.Sprintf("%s-%s-%s", categoryParent.Slug, product.Brand.Slug, *dto.Title))

		slug = &newSlug
	}

	ex = ps.repo.UpdateProduct(dto, slug, product.Id)

	if ex != nil {
		return ex
	}

	return nil
}

func (ps *ProductService) UpdateProductModel(dto UpdateProductModelDto, modelId int) exception.Error {

	model, err := ps.repo.FindModelById(modelId)

	if err != nil {
		return err
	}

	err = ps.repo.UpdateProductModel(dto, model.Id)

	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductService) DeleteProduct(id int) exception.Error {
	ex := pr.repo.DeleteProduct(id)
	if ex != nil {
		return ex
	}
	return nil
}

func (pr *ProductService) DeleteProductModel(id int) exception.Error {
	ex := pr.repo.DeleteProductModel(id)
	if ex != nil {
		return ex
	}
	return nil
}
