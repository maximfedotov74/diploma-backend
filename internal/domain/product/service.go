package product

import (
	"fmt"

	"github.com/maximfedotov74/fiber-psql/internal/domain/brand"
	"github.com/maximfedotov74/fiber-psql/internal/domain/category"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Repository interface {
	CreateProduct(dto CreateProductDto, slug string) error
	FindByProductSlugAndModelId(slug string, modelId int) (*Product, error)
	CreateModel(dto CreateProductModelDto) error
	AddPhoto(dto CreateProducModelImg) error
	FindById(id int) (*ProductWithoutRelations, error)
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

	err := ps.repo.CreateProduct(dto, slug)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (ps *ProductService) FindByProductSlugAndModelId(slug string, id int) (*Product, exception.Error) {
	p, err := ps.repo.FindByProductSlugAndModelId(slug, id)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}
	return p, nil
}

func (ps *ProductService) FindById(id int) (*ProductWithoutRelations, exception.Error) {
	p, err := ps.repo.FindById(id)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}
	return p, nil
}

func (ps *ProductService) CreateModel(dto CreateProductModelDto) exception.Error {
	_, ex := ps.FindById(dto.ProductId)

	if ex != nil {
		return ex
	}

	err := ps.repo.CreateModel(dto)

	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (ps *ProductService) AddPhoto(dto CreateProducModelImg) exception.Error {
	err := ps.repo.AddPhoto(dto)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil

}
