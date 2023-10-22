package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
)

type User interface {
	GetAll() error
	Create(dto model.CreateUserDto) (*UserRepoResponse, error)
	GetUserById(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type Role interface {
	Create(dto model.CreateRoleDto) (*model.Role, error)
	AddRoleToUser(roleId int, userId int) (bool, error)
	RemoveRoleFromUser(roleId int, userId int) (bool, error)
	FindRoleByTitle(title string) (*model.Role, error)
}

type Repositories struct {
	UserRepository User
	RoleRepository Role
}

func New(db *pgxpool.Pool) *Repositories {

	role := NewRoleRepository(db)
	user := NewUserRepository(db)

	return &Repositories{
		UserRepository: user,
		RoleRepository: role,
	}
}
