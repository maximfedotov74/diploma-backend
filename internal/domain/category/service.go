package category

import (
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Repository interface {
	CreateCategory(dto CreateCategoryDto, slug string) error
	FindByField(field string, value any) (*Category, error)
	GetCatalogCategories() ([]CatalogCategory, error)
}

type CategoryService struct {
	repo Repository
}

func NewCategoryService(repo Repository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (cs *CategoryService) CreateCategory(dto CreateCategoryDto) exception.Error {

	slug := utils.GenerateSlug(dto.Title)

	err := cs.repo.CreateCategory(dto, slug)

	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (cs *CategoryService) GetCatalogCategories() ([]CatalogCategory, exception.Error) {
	cts, err := cs.repo.GetCatalogCategories()

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return cts, nil
}

func (cs *CategoryService) FindBySlug(slug string) (*Category, exception.Error) {

	result, err := cs.repo.FindByField("slug", slug)
	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}
	return result, nil
}

func (cs *CategoryService) FindById(id int) (*Category, exception.Error) {

	result, err := cs.repo.FindByField("category_id", id)
	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}
	return result, nil
}
