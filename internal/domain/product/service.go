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
	FindBySlug(slug string) (*Product, error)
}

type CategoryService interface {
	FindById(id int) (*category.Category, exception.Error)
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
	brand, ex := ps.brandService.FindById(dto.BrandID)
	if ex != nil {
		return ex
	}

	slug := utils.GenerateSlug(fmt.Sprintf("%s-%s-%s", category.Slug, brand.Slug, dto.Title))

	err := ps.repo.CreateProduct(dto, slug)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (ps *ProductService) FindBySlug(slug string) (*Product, exception.Error) {
	p, err := ps.repo.FindBySlug(slug)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return p, nil
}
