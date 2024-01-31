package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type optionRepository interface {
	GetAll(ctx context.Context) ([]model.Option, fall.Error)
	CreateOption(ctx context.Context, dto model.CreateOptionDto) fall.Error
	UpdateOption(ctx context.Context, dto model.UpdateOptionDto, id int) fall.Error
	UpdateOptionValue(ctx context.Context, dto model.UpdateOptionValueDto, id int) fall.Error
	CreateValue(ctx context.Context, dto model.CreateOptionValueDto) fall.Error
	CreateSize(ctx context.Context, value string) fall.Error
	DeleteOption(ctx context.Context, id int) fall.Error
	DeleteValue(ctx context.Context, id int) fall.Error
	DeleteSize(ctx context.Context, id int) fall.Error
	DeleteSizeFromProductModel(ctx context.Context, modelSizeId int) fall.Error
	DeleteOptionFromProductModel(ctx context.Context, productModelOptionId int) fall.Error
	FindByField(ctx context.Context, field string, value any) (*model.Option, fall.Error)
	AddOptionToProductModel(ctx context.Context, dto model.AddOptionToProductModelDto) fall.Error
	AddSizeToProductModel(ctx context.Context, dto model.AddSizeToProductModelDto) fall.Error
	GetCatalogFilters(ctx context.Context, categorySlug string) (*model.CatalogFilters, fall.Error)
	CheckValueInOption(ctx context.Context, valueId int, optionId int) fall.Error
}

type OptionService struct {
	repo optionRepository
}

func NewOptionService(repo optionRepository) *OptionService {
	return &OptionService{repo: repo}
}

func (s *OptionService) GetCatalogFilters(ctx context.Context, slug string) (*model.CatalogFilters, fall.Error) {
	return s.repo.GetCatalogFilters(ctx, slug)
}

func (s *OptionService) GetAll(ctx context.Context) ([]model.Option, fall.Error) {
	return s.repo.GetAll(ctx)
}

func (s *OptionService) FindOptionById(ctx context.Context, id int) (*model.Option, fall.Error) {
	return s.repo.FindByField(ctx, "option_id", id)
}

func (s *OptionService) FindOptionBySlug(ctx context.Context, slug string) (*model.Option, fall.Error) {
	return s.repo.FindByField(ctx, "slug", slug)
}

func (s *OptionService) CreateOption(ctx context.Context, dto model.CreateOptionDto) fall.Error {
	option, _ := s.FindOptionBySlug(ctx, dto.Slug)

	if option != nil {
		return fall.NewErr(msg.OptionAlreadyExists, fall.STATUS_BAD_REQUEST)
	}

	ex := s.repo.CreateOption(ctx, dto)

	return ex
}

func (s *OptionService) CreateSize(ctx context.Context, dto model.CreateSizeDto) fall.Error {
	return s.repo.CreateSize(ctx, dto.Value)
}

func (s *OptionService) CreateValue(ctx context.Context, dto model.CreateOptionValueDto) fall.Error {
	return s.repo.CreateValue(ctx, dto)
}

func (s *OptionService) DeleteOption(ctx context.Context, id int) fall.Error {
	return s.repo.DeleteOption(ctx, id)
}

func (s *OptionService) DeleteValue(ctx context.Context, id int) fall.Error {
	return s.repo.DeleteValue(ctx, id)
}

func (s *OptionService) DeleteSize(ctx context.Context, id int) fall.Error {
	return s.repo.DeleteSize(ctx, id)
}

func (s *OptionService) DeleteSizeFromProductModel(ctx context.Context, modelSizeId int) fall.Error {
	return s.repo.DeleteSizeFromProductModel(ctx, modelSizeId)
}

func (s *OptionService) DeleteOptionFromProductModel(ctx context.Context, productModelOptionId int) fall.Error {
	return s.repo.DeleteOptionFromProductModel(ctx, productModelOptionId)
}

func (s *OptionService) AddOptionToProductModel(ctx context.Context, dto model.AddOptionToProductModelDto) fall.Error {

	err := s.repo.CheckValueInOption(ctx, dto.ValueId, dto.OptionId)

	if err != nil {
		return err
	}

	return s.repo.AddOptionToProductModel(ctx, dto)
}

func (s *OptionService) AddSizeToProductModel(ctx context.Context, dto model.AddSizeToProductModelDto) fall.Error {
	return s.repo.AddSizeToProductModel(ctx, dto)
}

func (s *OptionService) UpdateOption(ctx context.Context, dto model.UpdateOptionDto, id int) fall.Error {

	if dto.Slug != nil {
		opt, _ := s.FindOptionBySlug(ctx, *dto.Slug)
		if opt != nil {
			return fall.NewErr(msg.OptionAlreadyExists, fall.STATUS_BAD_REQUEST)
		}
	}
	return s.repo.UpdateOption(ctx, dto, id)
}

func (s *OptionService) UpdateOptionValue(ctx context.Context, dto model.UpdateOptionValueDto, id int) fall.Error {
	return s.repo.UpdateOptionValue(ctx, dto, id)
}
