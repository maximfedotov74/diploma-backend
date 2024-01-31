package repository

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

const (
	addRoleToUser = "INSERT INTO public.user_role (user_id, role_id) VALUES ($1, $2);"
)

type RoleRepository struct {
	db db.PostgresClient
}

func NewRoleRepository(db db.PostgresClient) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Create(ctx context.Context, dto model.CreateRoleDto) (*model.Role, fall.Error) {
	query := "INSERT INTO public.role (title) VALUES ($1) RETURNING role_id, title;"
	row := r.db.QueryRow(ctx, query, dto.Title)
	role := model.Role{}
	err := row.Scan(&role.Id, &role.Title)

	if err != nil {
		return nil, fall.NewErr(msg.RoleCreateError, fall.STATUS_INTERNAL_ERROR)
	}
	return &role, nil
}

func (r *RoleRepository) FindRoleByTitle(ctx context.Context, title string) (*model.Role, fall.Error) {
	query := `
	SELECT r.role_id, r.title, u.user_id, u.email FROM public.role as r
	LEFT JOIN user_role as ur ON r.role_id = ur.role_id
	LEFT JOIN public.user as u ON u.user_id = ur.user_id
	WHERE r.title = $1;`
	rows, err := r.db.Query(ctx, query, title)
	defer rows.Close()

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	founded := false
	role := model.Role{}

	users := []model.RoleUser{}

	for rows.Next() {
		user := model.RoleUser{}
		err := rows.Scan(&role.Id, &role.Title, &user.Id, &user.Email)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		if user.Id != nil {
			users = append(users, user)
		}
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.RoleNotFound, fall.STATUS_NOT_FOUND)
	}

	role.Users = users

	return &role, nil

}

func (r *RoleRepository) Find(ctx context.Context) ([]model.Role, fall.Error) {

	query := `
	SELECT r.role_id, r.title,  u.user_id, u.email, ur.user_role_id, ur.role_id as ur_role_id FROM public.role as r
	LEFT JOIN user_role as ur ON r.role_id = ur.role_id
	LEFT JOIN public.user as u ON u.user_id = ur.user_id;`

	rows, err := r.db.Query(ctx, query)
	defer rows.Close()

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	founded := false
	rolesMap := make(map[int]model.Role)
	usersMap := make(map[int]model.RoleUser)

	var rolesOrder []int
	var usersOrder []int

	for rows.Next() {
		var role model.Role
		role.Users = []model.RoleUser{}
		var user model.RoleUser

		err := rows.Scan(&role.Id, &role.Title, &user.Id, &user.Email, &user.UserRoleId, &user.RoleId)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		if user.UserRoleId != nil {
			_, ok := usersMap[*user.UserRoleId]
			if !ok {
				usersMap[*user.UserRoleId] = user
				usersOrder = append(usersOrder, *user.UserRoleId)
			}
		}
		_, ok := rolesMap[role.Id]
		if !ok {
			rolesMap[role.Id] = role
			rolesOrder = append(rolesOrder, role.Id)
		}
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	roles := make([]model.Role, 0, len(rolesMap))

	if !founded {
		return roles, nil
	}

	for _, key := range usersOrder {
		if err == nil {
			user := usersMap[key]
			role := rolesMap[*user.RoleId]
			role.Users = append(role.Users, user)
			rolesMap[role.Id] = role
		}
	}

	for _, key := range rolesOrder {
		role := rolesMap[key]
		roles = append(roles, role)
	}

	return roles, nil

}

func (r *RoleRepository) AddRoleToUser(ctx context.Context, roleId int, userId int, tx db.Transaction) fall.Error {
	if tx != nil {
		_, err := tx.Exec(ctx, addRoleToUser, userId, roleId)
		if err != nil {
			return fall.NewErr(msg.RoleAddError, fall.STATUS_INTERNAL_ERROR)
		}
		return nil
	}
	_, err := r.db.Exec(ctx, addRoleToUser, userId, roleId)
	if err != nil {
		return fall.NewErr(msg.RoleAddError, fall.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (r *RoleRepository) RemoveRole(ctx context.Context, roleId int) fall.Error {
	query := "DELETE FROM public.role WHERE role_id = $1;"

	_, err := r.db.Exec(ctx, query, roleId)

	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *RoleRepository) RemoveRoleFromUser(ctx context.Context, roleId int, userId int) fall.Error {

	query := "DELETE FROM public.user_role WHERE user_id = $1 AND role_id = $2;"
	_, err := r.db.Exec(ctx, query, userId, roleId)

	if err != nil {
		return fall.NewErr(msg.RoleDeleteError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}
