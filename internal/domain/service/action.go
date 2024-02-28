package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type actionRepository interface {
	Create(ctx context.Context, dto model.CreateActionDto) fall.Error
	AddModel(ctx context.Context, actionId string, modelId int) fall.Error
	FindById(ctx context.Context, id string) (*model.Action, fall.Error)
	GetAll(ctx context.Context) ([]model.Action, fall.Error)
	Update(ctx context.Context, dto model.UpdateActionDto, id string) fall.Error
	GetModels(ctx context.Context, id string) ([]model.ActionModel, fall.Error)
}

type actionProductService interface {
	FindProductModelById(ctx context.Context, id int) (*model.ProductModel, fall.Error)
}

type ActionService struct {
	repo           actionRepository
	productService actionProductService
}

func NewActionService(repo actionRepository, productService actionProductService) *ActionService {
	return &ActionService{repo: repo, productService: productService}
}

func (s *ActionService) GetModels(ctx context.Context, id string) ([]model.ActionModel, fall.Error) {
	return s.repo.GetModels(ctx, id)
}

func (s *ActionService) Create(ctx context.Context, dto model.CreateActionDto) fall.Error {
	return s.repo.Create(ctx, dto)
}

func (s *ActionService) AddModel(ctx context.Context, dto model.AddModelToActionDto) fall.Error {
	model, ex := s.productService.FindProductModelById(ctx, dto.ProductModelId)
	if ex != nil {
		return ex
	}
	action, ex := s.repo.FindById(ctx, dto.ActionId)
	if ex != nil {
		return ex
	}
	return s.repo.AddModel(ctx, action.Id, model.Id)
}

func (s *ActionService) GetAll(ctx context.Context) ([]model.Action, fall.Error) {
	return s.repo.GetAll(ctx)
}

func (s *ActionService) Update(ctx context.Context, dto model.UpdateActionDto, id string) fall.Error {
	_, ex := s.repo.FindById(ctx, id)
	if ex != nil {
		return ex
	}

	return s.repo.Update(ctx, dto, id)

}
