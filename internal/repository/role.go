package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (rr *RoleRepository) Create(dto model.CreateRoleDto) (*model.Role, error) {
	query := "INSERT INTO public.role (title) VALUES ($1) RETURNING role_id, title;"
	row := rr.db.QueryRow(context.Background(), query, dto.Title)
	role := model.Role{}
	err := row.Scan(&role.Id, &role.Title)

	if err != nil {
		return nil, errors.New(messages.ROLE_CREATE_ERROR)
	}
	return &role, nil
}

func (rr *RoleRepository) FindRoleByTitle(title string) (*model.Role, error) {
	row := rr.db.QueryRow(context.Background(), findRoleByTitle, title)
	role := model.Role{}
	err := row.Scan(&role.Id, &role.Title)
	if err != nil {
		return nil, errors.New(messages.ROLE_NOT_FOUND)
	}

	return &role, nil

}

func (rr *RoleRepository) AddRoleToUser(roleId int, userId int) (bool, error) {
	_, err := rr.db.Exec(context.Background(), addRoleToUser, userId, roleId)
	if err != nil {
		return false, errors.New(messages.ROLE_ADD_ERROR)
	}

	return true, nil
}

func (rr *RoleRepository) RemoveRoleFromUser(roleId int, userId int) (bool, error) {

	query := "DELETE FROM public.user_role WHERE user_id = $1 AND role_id = $2;"
	_, err := rr.db.Exec(context.Background(), query, userId, roleId)

	if err != nil {
		return false, errors.New(messages.ROLE_DELETE_ERROR)
	}

	return true, nil
}
