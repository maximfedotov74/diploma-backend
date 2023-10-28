package service

import (
	"strings"

	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

type RoleService struct {
	repo repository.Role
}

func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (rs *RoleService) Create(dto model.CreateRoleDto) (*model.Role, lib.Error) {

	dto.Title = strings.ToUpper(dto.Title)

	oldRole, _ := rs.repo.FindRoleByTitle(dto.Title)

	if oldRole != nil {
		return nil, lib.NewErr(messages.ROLE_EXISTS, 400)
	}

	role, err := rs.repo.Create(dto)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	return role, nil
}

func (rs *RoleService) AddRoleToUser(title string, userId int) lib.Error {

	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return lib.NewErr(err.Error(), 404)
	}

	err = rs.repo.AddRoleToUser(role.Id, userId)

	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	return nil

}

func (rs *RoleService) RemoveRoleFromUser(title string, userId int) lib.Error {
	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return lib.NewErr(err.Error(), 404)
	}

	err = rs.repo.RemoveRoleFromUser(role.Id, userId)

	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	return nil
}
