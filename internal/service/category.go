package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (cs *CategoryService) FindTypeByTitle(title string) (*model.CategoryType, lib.Error) {
	categoryType, err := cs.repo.FindTypeByTitle(title)
	if err != nil {
		return nil, lib.NewErr(err.Error(), 404)
	}
	return categoryType, nil
}

func (cs *CategoryService) FindCategoryByTitle(title string) (*model.Category, lib.Error) {
	category, err := cs.repo.FindCategoryByTitle(title)
	if err != nil {
		return nil, lib.NewErr(err.Error(), 404)
	}
	return category, nil
}

func (cs *CategoryService) CreateCategoryType(dto model.CreateCategoryTypeDto) lib.Error {

	err := cs.repo.CreateCategoryType(dto.Title)
	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	return nil
}

func (cs *CategoryService) CreateCategory(dto model.CreateCategoryDto) lib.Error {

	err := cs.repo.CreateCategory(dto.Title, dto.ImgPath, dto.ParentId)

	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	return nil
}
