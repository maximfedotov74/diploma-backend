package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
)

type RoleService struct {
	repo repository.Role
}

func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (rs *RoleService) Create(dto model.CreateRoleDto) (*model.Role, error) {
	role, err := rs.repo.Create(dto)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (rs *RoleService) AddRoleToUser(title string, userId int) (bool, error) {

	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return false, err
	}

	flag, err := rs.repo.AddRoleToUser(role.Id, userId)

	if err != nil {
		return flag, err
	}

	return flag, nil

}

func (rs *RoleService) RemoveRoleFromUser(title string, userId int) (bool, error) {
	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return false, err
	}

	flag, err := rs.repo.RemoveRoleFromUser(role.Id, userId)

	if err != nil {
		return flag, err
	}

	return flag, nil
}
