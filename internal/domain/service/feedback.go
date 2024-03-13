package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type feedbackRepository interface {
	AddFeedback(ctx context.Context, userId int, dto model.AddFeedbackDto) fall.Error
	GetModelFeedback(ctx context.Context, modelId int, order string) (*model.ModelFeedbackResponse, fall.Error)
	GetAll(ctx context.Context, order string, page int, filter string) (*model.AdminAllFeedbackResponse, fall.Error)
	DeleteFeedback(ctx context.Context, feedbackId int) fall.Error
	ToggleHidden(ctx context.Context, feedbackId int) fall.Error
	FindFeedback(ctx context.Context, userId int, modelId int) (*model.Feedback, fall.Error)
}

type FeedbackService struct {
	repo feedbackRepository
}

func NewFeedbackService(repo feedbackRepository) *FeedbackService {
	return &FeedbackService{
		repo: repo,
	}
}

func (s *FeedbackService) AddFeedback(ctx context.Context, dto model.AddFeedbackDto, userId int) fall.Error {

	f, _ := s.repo.FindFeedback(ctx, userId, dto.ModelId)

	if f != nil {
		return fall.NewErr(msg.FeedbackAlreadyExist, fall.STATUS_BAD_REQUEST)
	}

	return s.repo.AddFeedback(ctx, userId, dto)

}

func (s *FeedbackService) ToggleHidden(ctx context.Context, feedbackId int) fall.Error {
	return s.repo.ToggleHidden(ctx, feedbackId)
}

func (s *FeedbackService) GetModelFeedback(ctx context.Context, modelId int, order string) (*model.ModelFeedbackResponse, fall.Error) {
	return s.repo.GetModelFeedback(ctx, modelId, order)
}

func (s *FeedbackService) GetAll(ctx context.Context, order string, page int, filter string) (*model.AdminAllFeedbackResponse, fall.Error) {
	return s.repo.GetAll(ctx, order, page, filter)
}

func (s *FeedbackService) DeleteFeedback(ctx context.Context, feedbackId int) fall.Error {
	return s.repo.DeleteFeedback(ctx, feedbackId)

}
