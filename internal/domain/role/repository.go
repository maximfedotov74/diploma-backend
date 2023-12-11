package role

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/shared/db"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

const (
	addRoleToUser = "INSERT INTO public.user_role (user_id, role_id) VALUES ($1, $2);"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (rr *RoleRepository) Create(dto CreateRoleDto) (*Role, exception.Error) {
	query := "INSERT INTO public.role (title) VALUES ($1) RETURNING role_id, title;"
	row := rr.db.QueryRow(context.Background(), query, dto.Title)
	role := Role{}
	err := row.Scan(&role.Id, &role.Title)

	if err != nil {
		return nil, exception.NewErr(roleCreateError, exception.STATUS_INTERNAL_ERROR)
	}
	return &role, nil
}

func (rr *RoleRepository) FindRoleByTitle(title string) (*Role, exception.Error) {

	query := "SELECT role_id, title FROM public.role WHERE public.role.title = $1;"

	row := rr.db.QueryRow(context.Background(), query, title)
	role := Role{}

	err := row.Scan(&role.Id, &role.Title)

	if err != nil {
		return nil, exception.NewErr(roleNotFound, exception.STATUS_NOT_FOUND)
	}

	return &role, nil

}

func (rr *RoleRepository) AddRoleToUser(roleId int, userId int, tx *db.Transaction) exception.Error {

	if tx != nil {
		_, err := tx.Executer.Exec(tx.Ctx, addRoleToUser, userId, roleId)
		if err != nil {
			return exception.NewErr(roleAddError, exception.STATUS_INTERNAL_ERROR)
		}
		return nil
	}

	_, err := rr.db.Exec(context.Background(), addRoleToUser, userId, roleId)
	if err != nil {
		return exception.NewErr(roleAddError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (rr *RoleRepository) RemoveRoleFromUser(roleId int, userId int) exception.Error {

	query := "DELETE FROM public.user_role WHERE user_id = $1 AND role_id = $2;"
	_, err := rr.db.Exec(context.Background(), query, userId, roleId)

	if err != nil {
		return exception.NewErr(roleDeleteError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}
