package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/keys"
)

type userRoleRepository interface {
	FindRoleByTitle(ctx context.Context, title string) (*model.Role, fall.Error)
	AddRoleToUser(ctx context.Context, roleId int, userId int, tx db.Transaction) fall.Error
}

type UserRepository struct {
	db       db.PostgresClient
	roleRepo userRoleRepository
}

func NewUserRepository(db db.PostgresClient, roleRepo userRoleRepository) *UserRepository {
	return &UserRepository{db: db, roleRepo: roleRepo}
}

func (r *UserRepository) Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error) {

	tx, err := r.db.Begin(ctx)

	var ex fall.Error = nil

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	defer func() {
		if ex != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	q := "INSERT INTO public.user (email, password_hash) VALUES ($1,$2) RETURNING user_id, email;"

	row := tx.QueryRow(ctx, q, dto.Email, dto.Password)

	var email string
	var id int
	err = row.Scan(&id, &email)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	role, ex := r.roleRepo.FindRoleByTitle(ctx, keys.USER_ROLE)

	if ex != nil {
		return nil, ex
	}

	ex = r.roleRepo.AddRoleToUser(ctx, role.Id, id, tx)

	if ex != nil {
		return nil, ex
	}
	q = "INSERT INTO user_activation (user_id, activation_account_link) VALUES ($1, uuid_generate_v4()) RETURNING activation_account_link;"

	row = tx.QueryRow(ctx, q, id)

	var link string

	err = row.Scan(&link)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	return &model.CreatedUserResponse{Id: id, Email: email, Link: link}, nil

}

func (r *UserRepository) Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error {

	var queries []string

	if dto.AvatarPath != nil {
		queries = append(queries, fmt.Sprintf("avatar_path = '%s'", *dto.AvatarPath))
	}

	if dto.Patronymic != nil {
		queries = append(queries, fmt.Sprintf("patronymic = '%s'", *dto.Patronymic))
	}

	if dto.LastName != nil {
		queries = append(queries, fmt.Sprintf("last_name = '%s'", *dto.LastName))
	}

	if dto.FirstName != nil {
		queries = append(queries, fmt.Sprintf("first_name = '%s'", *dto.FirstName))
	}

	if dto.Gender != nil {
		queries = append(queries, fmt.Sprintf("gender = '%s'", *dto.Gender))
	}

	if len(queries) > 0 {
		q := "UPDATE public.user SET " + strings.Join(queries, ",") + " WHERE user_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.BrandUpdateError, err.Error()))
		}
		return nil
	}

	return nil
}

func (ur *UserRepository) findByIdOrEmail(ctx context.Context, field string, value any) (*model.User, fall.Error) {

	query := fmt.Sprintf(`
	SELECT public.user.user_id, public.user.email, public.user.password_hash,
	public.user.patronymic, public.user.first_name,	public.user.last_name, 
	role.title, role.role_id, public.user.is_activated, public.user.gender, public.user.avatar_path
	FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	WHERE public.user.%s = $1;`, field)

	rows, err := ur.db.Query(ctx, query, value)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	user := model.User{}
	founded := false
	for rows.Next() {
		role := model.UserRole{}
		err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Patronymic, &user.FirstName, &user.LastName,
			&role.Title, &role.Id, &user.IsActivated, &user.Gender, &user.AvatarPath)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		if role.Id != nil {
			user.Roles = append(user.Roles, role)
		}
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.UserNotFound, fall.STATUS_NOT_FOUND)
	}

	return &user, nil

}

func (r *UserRepository) FindById(ctx context.Context, id int) (*model.User, fall.Error) {
	return r.findByIdOrEmail(ctx, "user_id", id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, fall.Error) {
	return r.findByIdOrEmail(ctx, "email", email)
}

func (r *UserRepository) RemoveChangePasswordCode(ctx context.Context, userId int, tx db.Transaction) error {
	query := "DELETE FROM change_password_code WHERE user_id = $1"

	if tx != nil {
		_, err := tx.Exec(ctx, query, userId)

		if err != nil {
			return err
		}
		return nil
	}

	_, err := r.db.Exec(context.Background(), query, userId)

	if err != nil {
		return err
	}
	return nil

}

func (r *UserRepository) FindChangePasswordCode(ctx context.Context, userId int, code string) (*model.ChangePasswordCode, fall.Error) {
	query := `SELECT change_password_code_id, code, user_id FROM
	change_password_code WHERE user_id = $1 AND code = $2 AND end_time > CURRENT_TIMESTAMP;`

	row := r.db.QueryRow(ctx, query, userId, code)

	codeModel := model.ChangePasswordCode{}

	err := row.Scan(&codeModel.ChangePasswordCodeId, &codeModel.Code, &codeModel.UserId)

	if err != nil {
		return nil, fall.NewErr(msg.ChangePasswordCodeNotFound, fall.STATUS_NOT_FOUND)
	}

	return &codeModel, nil

}

func (r *UserRepository) CreateChangePasswordCode(ctx context.Context, userId int) (*string, fall.Error) {

	var ex fall.Error = nil

	tx, err := r.db.Begin(ctx)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	defer func() {
		if ex != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = r.RemoveChangePasswordCode(ctx, userId, nil)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	query := "INSERT INTO change_password_code (user_id) VALUES ($1) RETURNING code;"

	row := tx.QueryRow(ctx, query, userId)

	var code string

	err = row.Scan(&code)

	if err != nil {
		ex = fall.ServerError(msg.CreateChangeCodeError)
		return nil, ex
	}

	return &code, nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, userId int, newPassword string) fall.Error {

	query := `UPDATE public.user SET password_hash = $1,
	updated_at = CURRENT_TIMESTAMP WHERE public.user.user_id = $2;`

	_, err := r.db.Exec(ctx, query, newPassword, userId)
	if err != nil {
		return fall.NewErr(msg.UpdatePasswordError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context, page int) (*model.GetAllUsersResponse, fall.Error) {
	limit := 16

	offset := page*limit - limit

	q := "SELECT user_id, (select COUNT(*) from public.user) as total FROM public.user ORDER BY public.user.created_at LIMIT $1 OFFSET $2;"

	rows, err := r.db.Query(ctx, q, limit, offset)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var userIds []int
	var total int

	for rows.Next() {
		var id int
		err := rows.Scan(&id, &total)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		userIds = append(userIds, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	q = `SELECT public.user.user_id, public.user.email, public.user.password_hash,
	public.user.patronymic, public.user.first_name,	public.user.last_name, 
	role.title, role.role_id, public.user.is_activated, public.user.gender, public.user.avatar_path, user_role.user_id as role_user_id,
	user_role.user_role_id as user_role_id
	FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	WHERE public.user.user_id = ANY ($1);`

	rows, err = r.db.Query(ctx, q, userIds)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	rolesMap := make(map[int]model.UserRole)
	var rolesOrder []int
	usersMap := make(map[int]*model.User)

	for rows.Next() {
		role := model.UserRole{}
		user := model.User{}

		err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Patronymic, &user.FirstName, &user.LastName,
			&role.Title, &role.Id, &user.IsActivated, &user.Gender, &user.AvatarPath, &role.UserId, &role.UserRoleId)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		_, ok := usersMap[user.Id]
		if !ok {
			usersMap[user.Id] = &user
		}

		if role.Id != nil && role.UserId != nil && role.UserRoleId != nil {
			_, ok := rolesMap[*role.UserRoleId]
			if !ok {
				rolesMap[*role.UserRoleId] = role
				rolesOrder = append(rolesOrder, *role.UserRoleId)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	for _, key := range rolesOrder {
		r := rolesMap[key]
		u := usersMap[*r.UserId]
		u.Roles = append(u.Roles, r)
	}

	result := make([]*model.User, 0, len(usersMap))

	for _, key := range userIds {
		u := usersMap[key]
		result = append(result, u)
	}

	return &model.GetAllUsersResponse{
		Users: result,
		Total: total,
	}, nil

}
