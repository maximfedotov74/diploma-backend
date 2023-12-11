package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/domain/role"
	"github.com/maximfedotov74/fiber-psql/internal/shared/db"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

// TODO ADD erros and update find method
const (
	userRole  = "USER"
	adminRole = "ADMIN"
)

const (
	emailField  = "email"
	userIdField = "user_id"
)

type RoleRepository interface {
	FindRoleByTitle(title string) (*role.Role, exception.Error)
	AddRoleToUser(roleId int, userId int, tx *db.Transaction) exception.Error
}

type UserRepository struct {
	db       *pgxpool.Pool
	roleRepo RoleRepository
}

func NewUserRepository(db *pgxpool.Pool, roleRepo RoleRepository) *UserRepository {
	return &UserRepository{db: db, roleRepo: roleRepo}
}

func (ur *UserRepository) GetAll() error {
	return nil
}

func (ur *UserRepository) Create(password string, email string) (*UserCreatedResponse, exception.Error) {
	txCtx := context.Background()
	commit := false

	tx, err := ur.db.Begin(txCtx)

	defer func() {
		if !commit {
			tx.Rollback(txCtx)
		} else {
			tx.Commit(txCtx)
		}
	}()

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	query := "INSERT INTO public.user (email, password_hash) VALUES ($1, $2) RETURNING user_id, email;"

	row := tx.QueryRow(txCtx, query, email, password)
	var id int
	var userEmail string

	err = row.Scan(&id, &userEmail)
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	role, ex := ur.roleRepo.FindRoleByTitle(userRole)

	if ex != nil {
		return nil, ex
	}

	ex = ur.roleRepo.AddRoleToUser(role.Id, id, &db.Transaction{Executer: tx, Ctx: txCtx})

	if err != nil {
		return nil, ex
	}

	query = "INSERT INTO user_settings (user_id) VALUES ($1);"

	_, err = tx.Exec(txCtx, query, id)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	query = "INSERT INTO public.user_activation (user_id, activation_account_link) VALUES ($1, uuid_generate_v4()) RETURNING activation_account_link;"
	row = tx.QueryRow(txCtx, query, id)
	var link string
	err = row.Scan(&link)
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	commit = true
	return &UserCreatedResponse{Id: id, ActivationAccountLink: link, Email: userEmail}, nil
}

func (ur *UserRepository) findByIdOrEmail(field string, value any) (*User, exception.Error) {

	ctx := context.Background()

	query := fmt.Sprintf(`SELECT public.user.user_id, public.user.email, public.user.password_hash,
	role.title, role.role_id, public.user_settings.is_activated
	FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	LEFT JOIN public.user_settings ON public.user.user_id = public.user_settings.user_id
	WHERE public.user.%s = $1;`, field)

	rows, err := ur.db.Query(ctx, query, value)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	defer rows.Close()

	user := User{}
	founded := false
	for rows.Next() {
		role := role.Role{}
		err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &role.Title, &role.Id, &user.IsActivated)
		if err != nil {
			return nil, exception.ServerError(err.Error())
		}
		user.Roles = append(user.Roles, role)
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}
	if !founded {
		return nil, exception.NewErr(userNotFound, exception.STATUS_NOT_FOUND)
	}

	return &user, nil

}

func (ur *UserRepository) GetUserById(id int) (*User, exception.Error) {
	user, err := ur.findByIdOrEmail(userIdField, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*User, exception.Error) {
	user, err := ur.findByIdOrEmail(emailField, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FindActivationLink(link string) (*int, exception.Error) {
	ctx := context.Background()
	query := `SELECT user_activation.user_id FROM user_activation
	WHERE user_activation.activation_account_link = $1;`
	row := ur.db.QueryRow(ctx, query, link)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, exception.NewErr(activationNotFound, exception.STATUS_NOT_FOUND)
	}

	return &id, nil
}

func (ur *UserRepository) ActivateUser(id *int) exception.Error {

	txCtx := context.Background()

	tx, err := ur.db.Begin(context.Background())
	commit := false

	if err != nil {
		return exception.ServerError(err.Error())
	}

	defer func() {
		if !commit {
			tx.Rollback(txCtx)
		} else {
			tx.Commit(txCtx)
		}
	}()

	query := `UPDATE user_settings SET is_activated = TRUE WHERE user_id = $1;`

	_, err = tx.Exec(txCtx, query, id)

	if err != nil {
		return exception.NewErr(activationError, exception.STATUS_INTERNAL_ERROR)
	}

	query = "DELETE FROM user_activation WHERE user_id = $1;"
	_, err = tx.Exec(txCtx, query, id)

	if err != nil {
		return exception.NewErr(activationError, exception.STATUS_INTERNAL_ERROR)
	}
	commit = true
	return nil
}

func (ur *UserRepository) ChangePassword(userId int, newPassword string) exception.Error {

	txCtx := context.Background()

	tx, err := ur.db.Begin(txCtx)

	defer func() {
		if err != nil {
			tx.Rollback(txCtx)
		} else {
			tx.Commit(txCtx)
		}
	}()

	if err != nil {
		return exception.ServerError(err.Error())
	}

	query := `UPDATE public.user SET password_hash = $1,
	updated_at = CURRENT_TIMESTAMP WHERE public.user.user_id = $2;`

	_, err = tx.Exec(txCtx, query, newPassword, userId)
	if err != nil {
		return exception.NewErr(updatePasswordError, exception.STATUS_INTERNAL_ERROR)
	}

	err = removeChangePasswordCode(userId, tx, txCtx)
	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil
}

func removeChangePasswordCode(userId int, tx pgx.Tx, ctx context.Context) error {
	query := "DELETE FROM change_password_code WHERE user_id = $1"

	_, err := tx.Exec(ctx, query, userId)

	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindChangePasswordCode(userId int, code string) (*ChangePasswordCode, exception.Error) {
	ctx := context.Background()
	query := `SELECT change_password_code_id, code, user_id FROM
	change_password_code WHERE user_id = $1 AND code = $2 AND end_time > CURRENT_TIMESTAMP;`

	row := ur.db.QueryRow(ctx, query, userId, code)

	codeModel := ChangePasswordCode{}

	err := row.Scan(&codeModel.ChangePasswordCodeId, &codeModel.Code, &codeModel.UserId)

	if err != nil {
		return nil, exception.NewErr(changePasswordCodeNotFound, exception.STATUS_NOT_FOUND)
	}

	return &codeModel, nil

}

func (ur *UserRepository) CreateChangePasswordCode(userId int) (*string, exception.Error) {

	txCtx := context.Background()

	tx, err := ur.db.Begin(txCtx)

	defer func() {
		if err != nil {
			tx.Rollback(txCtx)
		} else {
			tx.Commit(txCtx)
		}
	}()

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	err = removeChangePasswordCode(userId, tx, txCtx)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	query := "INSERT INTO change_password_code (user_id) VALUES ($1) RETURNING code;"

	row := tx.QueryRow(txCtx, query, userId)

	var code string

	err = row.Scan(&code)

	if err != nil {
		return nil, exception.NewErr(createChangeCodeError, exception.STATUS_INTERNAL_ERROR)
	}

	return &code, nil
}
