package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type brandRepository interface {
	CreateBrand(ctx context.Context, dto model.CreateBrandDto, slug string) fall.Error
	FindByFeild(ctx context.Context, field string, value any) (*model.Brand, fall.Error)
	GetAll(ctx context.Context) ([]model.Brand, fall.Error)
	UpdateBrand(ctx context.Context, dto model.UpdateBrandDto, newSlug *string, id int) fall.Error
	Delete(ctx context.Context, slug string) fall.Error
}

type BrandService struct {
	repo brandRepository
}

func NewBrandService(repo brandRepository) *BrandService {
	return &BrandService{repo: repo}
}

func (s *BrandService) Create(ctx context.Context, dto model.CreateBrandDto) fall.Error {
	b, _ := s.FindByTitle(ctx, dto.Title)

	if b != nil {
		return fall.NewErr(msg.BrandExists, fall.STATUS_BAD_REQUEST)
	}

	slug := utils.GenerateSlug(dto.Title)

	err := s.repo.CreateBrand(ctx, dto, slug)

	return err

}

func (s *BrandService) Update(ctx context.Context, dto model.UpdateBrandDto, id int) fall.Error {
	_, err := s.FindById(ctx, id)

	if err != nil {
		return err
	}

	var slug *string

	if dto.Title != nil {
		b, _ := s.FindByTitle(ctx, *dto.Title)
		if b != nil {
			return fall.NewErr(msg.BrandExists, fall.STATUS_BAD_REQUEST)
		}
		newSlug := utils.GenerateSlug(*dto.Title)
		slug = &newSlug
	}

	err = s.repo.UpdateBrand(ctx, dto, slug, id)
	return err
}

func (s *BrandService) GetAll(ctx context.Context) ([]model.Brand, fall.Error) {
	return s.repo.GetAll(ctx)
}

func (s *BrandService) FindByTitle(ctx context.Context, title string) (*model.Brand, fall.Error) {
	return s.repo.FindByFeild(ctx, "title", title)
}

func (s *BrandService) FindBySlug(ctx context.Context, slug string) (*model.Brand, fall.Error) {
	return s.repo.FindByFeild(ctx, "slug", slug)
}
func (s *BrandService) FindById(ctx context.Context, id int) (*model.Brand, fall.Error) {
	return s.repo.FindByFeild(ctx, "brand_id", id)
}

func (s *BrandService) Delete(ctx context.Context, slug string) fall.Error {
	return s.repo.Delete(ctx, slug)
}
