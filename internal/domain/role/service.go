package role

import (
	"strings"

	"github.com/maximfedotov74/fiber-psql/internal/shared/db"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type Repository interface {
	Create(dto CreateRoleDto) (*Role, error)
	FindRoleByTitle(title string) (*Role, error)
	AddRoleToUser(roleId int, userId int, transaction *db.Transaction) error
	RemoveRoleFromUser(roleId int, userId int) error
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
		return nil, exception.NewErr(messages.ROLE_EXISTS, 400)
	}

	role, err := rs.repo.Create(dto)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return role, nil
}

func (rs *RoleService) AddRoleToUser(title string, userId int) exception.Error {

	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return exception.NewErr(err.Error(), 404)
	}

	err = rs.repo.AddRoleToUser(role.Id, userId, nil)

	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil

}

func (rs *RoleService) RemoveRoleFromUser(title string, userId int) exception.Error {
	role, err := rs.repo.FindRoleByTitle(title)

	if err != nil {
		return exception.NewErr(err.Error(), 404)
	}

	err = rs.repo.RemoveRoleFromUser(role.Id, userId)

	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}
