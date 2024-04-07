package service

import (
	"context"
	"strings"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type roleRepository interface {
	Create(ctx context.Context, dto model.CreateRoleDto) (*model.Role, fall.Error)
	FindRoleByTitle(ctx context.Context, title string) (*model.Role, fall.Error)
	AddRoleToUser(ctx context.Context, roleId int, userId int, tx db.Transaction) fall.Error
	RemoveRoleFromUser(ctx context.Context, roleId int, userId int) fall.Error
	FindWithUsers(ctx context.Context) ([]model.Role, fall.Error)
	RemoveRole(ctx context.Context, roleId int) fall.Error
	GetAll(ctx context.Context) ([]model.UserRole, fall.Error)
}

type RoleService struct {
	repo roleRepository
}

func NewRoleService(repo roleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAll(ctx context.Context) ([]model.UserRole, fall.Error) {
	return s.repo.GetAll(ctx)
}

func (s *RoleService) Create(ctx context.Context, dto model.CreateRoleDto) (*model.Role, fall.Error) {
	dto.Title = strings.ToUpper(dto.Title)

	existRole, _ := s.repo.FindRoleByTitle(ctx, dto.Title)

	if existRole != nil {
		return nil, fall.NewErr(msg.RoleExists, fall.STATUS_BAD_REQUEST)
	}

	newRole, err := s.repo.Create(ctx, dto)

	if err != nil {
		return nil, err
	}

	return newRole, nil
}

func (s *RoleService) FindRoleByTitle(ctx context.Context, title string) (*model.Role, fall.Error) {

	title = strings.ToUpper(title)

	role, err := s.repo.FindRoleByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) FindWithUsers(ctx context.Context) ([]model.Role, fall.Error) {
	return s.repo.FindWithUsers(ctx)
}

func (s *RoleService) RemoveRoleFromUser(ctx context.Context, title string, userId int) fall.Error {
	role, err := s.FindRoleByTitle(ctx, title)
	if err != nil {
		return err
	}
	err = s.repo.RemoveRoleFromUser(ctx, role.Id, userId)
	return err
}

func (s *RoleService) AddRoleToUser(ctx context.Context, title string, userId int) fall.Error {
	role, err := s.FindRoleByTitle(ctx, title)

	if err != nil {
		return err
	}

	err = s.repo.AddRoleToUser(ctx, role.Id, userId, nil)

	return err
}

func (s *RoleService) RemoveRole(ctx context.Context, roleId int) fall.Error {
	return s.repo.RemoveRole(ctx, roleId)
}
