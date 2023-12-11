package feedback

import exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"

type Repository interface {
	AddFeedback(userId int, dto AddFeedbackDto) exception.Error
	GetModelFeedback(modelId int, order string) (*ModelFeedbackResponse, exception.Error)
	GetAll(order string) ([]Feedback, exception.Error)
	DeleteFeedback(feedbackId int) exception.Error
	ToggleHidden(feedbackId int) exception.Error
}

type FeedbackService struct {
	repo Repository
}

func NewFeedbackService(repo Repository) *FeedbackService {
	return &FeedbackService{
		repo: repo,
	}
}

func (fs *FeedbackService) AddFeedback(dto AddFeedbackDto, userId int) exception.Error {
	ex := fs.repo.AddFeedback(userId, dto)
	if ex != nil {
		return ex
	}
	return nil
}

func (fs *FeedbackService) ToggleHidden(feedbackId int) exception.Error {
	ex := fs.repo.ToggleHidden(feedbackId)
	if ex != nil {
		return ex
	}
	return nil
}

func (fs *FeedbackService) GetModelFeedback(modelId int, order string) (*ModelFeedbackResponse, exception.Error) {
	return fs.repo.GetModelFeedback(modelId, order)
}

func (fs *FeedbackService) GetAll(order string) ([]Feedback, exception.Error) {
	return fs.repo.GetAll(order)
}

func (fs *FeedbackService) DeleteFeedback(feedbackId int) exception.Error {
	ex := fs.repo.DeleteFeedback(feedbackId)

	if ex != nil {
		return ex
	}

	return nil
}
