package role

import (
	"strings"

	"github.com/maximfedotov74/fiber-psql/internal/shared/db"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type Repository interface {
	Create(dto CreateRoleDto) (*Role, exception.Error)
	FindRoleByTitle(title string) (*Role, exception.Error)
	AddRoleToUser(roleId int, userId int, transaction *db.Transaction) exception.Error
	RemoveRoleFromUser(roleId int, userId int) exception.Error
}

type RoleService struct {
	repo Repository
}

func NewRoleService(repo Repository) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (rs *RoleService) Create(dto CreateRoleDto) (*Role, exception.Error) {

	dto.Title = strings.ToUpper(dto.Title)

	oldRole, _ := rs.repo.FindRoleByTitle(dto.Title)

	if oldRole != nil {
		return nil, exception.NewErr(roleExists, exception.STATUS_BAD_REQUEST)
	}

	role, err := rs.repo.Create(dto)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (rs *RoleService) AddRoleToUser(title string, userId int) exception.Error {

	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return err
	}

	err = rs.repo.AddRoleToUser(role.Id, userId, nil)

	if err != nil {
		return err
	}

	return nil

}

func (rs *RoleService) RemoveRoleFromUser(title string, userId int) exception.Error {
	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return err
	}

	err = rs.repo.RemoveRoleFromUser(role.Id, userId)

	if err != nil {
		return err
	}

	return nil
}
