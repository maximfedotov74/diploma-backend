package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type ActionRepository struct {
	db db.PostgresClient
}

func NewActionRepository(db db.PostgresClient) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) Create(ctx context.Context, dto model.CreateActionDto) fall.Error {

	q := "INSERT INTO action (end_date,title,img_path,description,action_gender) VALUES ($1,$2,$3,$4,$5);"

	_, err := r.db.Exec(ctx, q, dto.EndDate, dto.Title, dto.ImgPath, dto.Description, dto.Gender)
	if err != nil {
		return fall.ServerError(fmt.Sprintf("Ошибка при создании акции: %s", err.Error()))
	}

	return nil
}

func (r *ActionRepository) AddModel(ctx context.Context, actionId string, modelId int) fall.Error {
	q := "INSERT INTO action_model (action_id,product_model_id) VALUES ($1,$2);"

	_, err := r.db.Exec(ctx, q, actionId, modelId)
	if err != nil {
		return fall.ServerError(fmt.Sprintf("Ошибка при добавлении модели в акцию: %s", err.Error()))
	}
	return nil
}

func (r *ActionRepository) GetModels(ctx context.Context, id string) ([]model.ActionModel, fall.Error) {
	q := `
	SELECT p.product_id as p_id, p.title as p_title,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, ct.category_id as ct_id, ct.title as ct_title, ct.slug as ct_slug,
	pm.product_model_id as model_id, pm.slug as m_slug, pm.article as m_article, pm.price as model_price, pm.discount as model_discount,
	pm.main_image_path as pm_main_img, am.action_model_id as am_id
	FROM product p INNER JOIN category ct ON p.category_id = ct.category_id 
	INNER JOIN brand b on p.brand_id = b.brand_id
	INNER JOIN product_model pm ON pm.product_id = p.product_id
	inner join action_model as am on am.product_model_id = pm.product_model_id
	inner join action as a on a.action_id = am.action_id
	where a.action_id = $1 ORDER BY am.action_model_id;
	`

	rows, err := r.db.Query(ctx, q, id)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var models []model.ActionModel

	for rows.Next() {
		m := model.ActionModel{}

		err := rows.Scan(&m.ProductId, &m.Title, &m.Brand.Id, &m.Brand.Title, &m.Brand.Slug,
			&m.Category.Id, &m.Category.Title, &m.Category.Slug, &m.ModelId, &m.Slug, &m.Article, &m.Price, &m.Discount,
			&m.MainImagePath, &m.ActionModelId,
		)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		models = append(models, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return models, nil

}

func (r *ActionRepository) FindById(ctx context.Context, id string) (*model.Action, fall.Error) {
	q := `SELECT action_id, created_at, updated_at, end_date, title, is_activated, img_path, description FROM action WHERE action_id = $1;`

	row := r.db.QueryRow(ctx, q, id)

	action := model.Action{}

	err := row.Scan(&action.Id, &action.CreatedAt, &action.UpdatedAt, &action.EndDate, &action.Title, &action.IsActivated,
		&action.ImgPath, &action.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.ActionNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &action, nil
}

func (r *ActionRepository) GetAll(ctx context.Context) ([]model.Action, fall.Error) {
	q := `SELECT action_id, created_at, updated_at, end_date, title, is_activated, img_path, description, action_gender
	FROM action ORDER BY created_at;`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var actions []model.Action

	for rows.Next() {
		action := model.Action{}

		err := rows.Scan(&action.Id, &action.CreatedAt, &action.UpdatedAt, &action.EndDate, &action.Title, &action.IsActivated,
			&action.ImgPath, &action.Description, &action.Gender)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return actions, nil
}

func (r *ActionRepository) Update(ctx context.Context, dto model.UpdateActionDto, id string) fall.Error {
	var queries []string

	if dto.Description != nil {
		queries = append(queries, fmt.Sprintf("description = '%s'", *dto.Description))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = '%s'", *dto.Title))
	}

	if dto.ImgPath != nil {
		queries = append(queries, fmt.Sprintf("img_path = '%s'", *dto.ImgPath))
	}

	if dto.EndDate != nil {

		f := dto.EndDate.Format("2006-01-02 15:04:05")

		queries = append(queries, fmt.Sprintf("end_date = '%s'", f))
	}

	if dto.IsActivated != nil {
		queries = append(queries, fmt.Sprintf("is_activated = '%t'", *dto.IsActivated))
	}

	if len(queries) > 0 {
		q := "UPDATE action SET " + strings.Join(queries, ",") + " WHERE action_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.ActionUpdateError, err.Error()))
		}
	}
	return nil
}

func (r *ActionRepository) DeleteActionModel(ctx context.Context, actionModelId int) fall.Error {

	q := "DELETE FROM action_model WHERE action_model_id = $1;"

	_, err := r.db.Exec(ctx, q, actionModelId)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil

}

func (r *ActionRepository) DeleteAction(ctx context.Context, id string) fall.Error {
	q := "DELETE FROM action WHERE action_id = $1;"

	_, err := r.db.Exec(ctx, q, id)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *ActionRepository) GetActionsByGender(ctx context.Context, gender model.ActionGender) ([]model.Action, fall.Error) {

	q := `SELECT action_id, created_at, updated_at, end_date, title, is_activated, img_path, description, action_gender
	FROM action WHERE action_gender IN ($1, 'everyone')
	AND is_activated = TRUE AND current_timestamp < end_date
	ORDER BY action_gender;`

	rows, err := r.db.Query(ctx, q, gender)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var actions []model.Action

	for rows.Next() {
		action := model.Action{}

		err := rows.Scan(&action.Id, &action.CreatedAt, &action.UpdatedAt, &action.EndDate, &action.Title, &action.IsActivated,
			&action.ImgPath, &action.Description, &action.Gender)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return actions, nil
}
