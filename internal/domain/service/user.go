package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type userRepository interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
	FindByEmail(ctx context.Context, email string) (*model.User, fall.Error)
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
	Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error) {
	return s.repo.Create(ctx, dto)
}

func (s *UserService) FindById(ctx context.Context, id int) (*model.User, fall.Error) {
	return s.repo.FindById(ctx, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, fall.Error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error {
	return s.repo.Update(ctx, dto, id)
}
