package category

import (
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Repository interface {
	CreateCategory(dto CreateCategoryDto, slug string) exception.Error
	FindByFieldWithSubcategories(field string, value any) (*Category, exception.Error)
	FindByField(field string, value any) (*CategoryDb, exception.Error)
	GetParentTopLevel(parentId int) (*CategoryDb, exception.Error)
	GetParentSubLevel(id int) (*CategoryDb, exception.Error)
	GetAll() ([]Category, exception.Error)
	UpdateCategory(dto UpdateCategoryDto, newSlug *string, id int) exception.Error
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

	exist, _ := cs.FindByTitle(dto.Title)

	if exist != nil {
		return exception.NewErr(categoryTitleUnique, 400)
	}

	slug := utils.GenerateSlug(dto.Title)

	err := cs.repo.CreateCategory(dto, slug)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CategoryService) FindBySlugWithSubcategories(slug string) (*Category, exception.Error) {

	result, err := cs.repo.FindByFieldWithSubcategories("slug", slug)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CategoryService) FindByIdWithSubcategories(id int) (*Category, exception.Error) {

	result, err := cs.repo.FindByFieldWithSubcategories("category_id", id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CategoryService) FindBySlug(slug string) (*CategoryDb, exception.Error) {

	result, err := cs.repo.FindByField("slug", slug)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CategoryService) FindByTitle(title string) (*CategoryDb, exception.Error) {

	result, err := cs.repo.FindByField("title", title)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CategoryService) FindById(id int) (*CategoryDb, exception.Error) {

	result, err := cs.repo.FindByField("category_id", id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CategoryService) GetParentSubLevel(id int) (*CategoryDb, exception.Error) {
	cat, err := cs.repo.GetParentSubLevel(id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (cs *CategoryService) GetAll() ([]Category, exception.Error) {
	cats, err := cs.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (cs *CategoryService) UpdateCategory(dto UpdateCategoryDto, id int) exception.Error {

	_, err := cs.FindById(id)

	if err != nil {
		return err
	}

	var slug *string

	if dto.Title != nil {

		exist, _ := cs.FindByTitle(*dto.Title)

		if exist != nil {
			return exception.NewErr(categoryTitleUnique, 400)
		}

		newSlug := utils.GenerateSlug(*dto.Title)
		slug = &newSlug
	}

	err = cs.repo.UpdateCategory(dto, slug, id)

	if err != nil {
		return err
	}
	return nil

}
