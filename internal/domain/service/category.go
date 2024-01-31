package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type categoryRepository interface {
	Create(ctx context.Context, dto model.CreateCategoryDto, slug string) fall.Error
	FindByField(ctx context.Context, field string, value any) (*model.CategoryModel, fall.Error)
	FindByFieldRelation(ctx context.Context, field string, value any) (*model.Category, fall.Error)
	GetParentSubLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error)
	GetParentTopLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error)
	Update(ctx context.Context, dto model.UpdateCategoryDto, newSlug *string, id int) fall.Error
	GetAll(ctx context.Context) ([]*model.Category, fall.Error)
	GetCatalogCategories(ctx context.Context, id int, activeSlug string) (*model.СatalogCategory, fall.Error)
	Delete(ctx context.Context, slug string) fall.Error
}

type CategoryService struct {
	repo categoryRepository
}

func NewCategoryService(repo categoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, dto model.CreateCategoryDto) fall.Error {
	c, _ := s.FindByTitle(ctx, dto.Title)
	if c != nil {
		return fall.NewErr(msg.CategoryExists, fall.STATUS_BAD_REQUEST)
	}
	slug := utils.GenerateSlug(dto.Title)
	err := s.repo.Create(ctx, dto, slug)
	return err
}

func (s *CategoryService) Update(ctx context.Context, dto model.UpdateCategoryDto, id int) fall.Error {
	_, err := s.FindById(ctx, id)

	if err != nil {
		return err
	}

	var slug *string

	if dto.Title != nil {
		exist, _ := s.FindByTitle(ctx, *dto.Title)
		if exist != nil {
			return fall.NewErr(msg.CategoryTitleUnique, fall.STATUS_BAD_REQUEST)
		}
		newSlug := utils.GenerateSlug(*dto.Title)
		slug = &newSlug
	}
	err = s.repo.Update(ctx, dto, slug, id)
	return err
}

func (s *CategoryService) FindByTitle(ctx context.Context, title string) (*model.CategoryModel, fall.Error) {
	return s.repo.FindByField(ctx, "title", title)
}

func (s *CategoryService) FindById(ctx context.Context, id int) (*model.CategoryModel, fall.Error) {
	return s.repo.FindByField(ctx, "category_id", id)
}

func (s *CategoryService) FindBySlug(ctx context.Context, slug string) (*model.CategoryModel, fall.Error) {
	return s.repo.FindByField(ctx, "slug", slug)
}

func (s *CategoryService) GetParentSubLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error) {
	return s.repo.GetParentSubLevel(ctx, id)
}

func (s *CategoryService) GetParentTopLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error) {
	return s.repo.GetParentTopLevel(ctx, id)
}

func (s *CategoryService) FindByIdRelation(ctx context.Context, id int) (*model.Category, fall.Error) {
	return s.repo.FindByFieldRelation(ctx, "category_id", id)
}

func (s *CategoryService) FindBySlugRelation(ctx context.Context, slug string) (*model.Category, fall.Error) {
	return s.repo.FindByFieldRelation(ctx, "slug", slug)
}

func (s *CategoryService) Delete(ctx context.Context, slug string) fall.Error {
	return s.repo.Delete(ctx, slug)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*model.Category, fall.Error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) GetCatalogCategories(ctx context.Context, slug string) (*model.СatalogCategory, fall.Error) {
	current, ex := s.FindBySlug(ctx, slug)

	if ex != nil {
		return nil, ex
	}

	if current.ParentId != nil {
		parent, ex := s.repo.GetParentTopLevel(ctx, *current.ParentId)
		if ex != nil {
			return nil, ex
		}
		cat, ex := s.repo.GetCatalogCategories(ctx, parent.Id, current.Slug)
		if ex != nil {
			return nil, ex
		}
		return cat, nil
	} else {
		cat, ex := s.repo.GetCatalogCategories(ctx, current.Id, current.Slug)
		if ex != nil {
			return nil, ex
		}
		return cat, nil
	}
}
