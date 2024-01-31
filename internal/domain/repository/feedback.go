package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type FeedbackRepository struct {
	db db.PostgresClient
}

func NewFeedbackRepository(db db.PostgresClient) *FeedbackRepository {
	return &FeedbackRepository{
		db: db,
	}
}

func (r *FeedbackRepository) FindFeedback(ctx context.Context, userId int, modelId int) (*model.Feedback, fall.Error) {
	q := `SELECT f.feedback_id as f_id, f.created_at as created_at, f.updated_at as updated_at, f.feedback_text as f_text,
	f.rate as f_rate, f.product_model_id as f_model_id, f.is_hidden as f_hidden,
	u.user_id as u_id, u.email as u_email,
	u.avatar_path as u_avatar_path, u.first_name as u_first_name, u.last_name as u_last_name FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	WHERE
  f.user_id = $1 AND f.product_model_id = $2;
  `

	row := r.db.QueryRow(ctx, q, userId, modelId)

	f := model.Feedback{}

	err := row.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &f.ModelId, &f.Hidden,
		&f.User.Id, &f.User.Email, &f.User.Avatar, &f.User.FirstName, &f.User.LastName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.FeedbackNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &f, nil
}

func (r *FeedbackRepository) AddFeedback(ctx context.Context, userId int, dto model.AddFeedbackDto) fall.Error {
	query := `
	INSERT INTO feedback (feedback_text, rate, product_model_id, user_id) VALUES ($1,$2,$3,$4);
	`
	_, err := r.db.Exec(ctx, query, dto.Text, dto.Rate, dto.ModelId, userId)

	if err != nil {
		return fall.NewErr(msg.FeedbackCreateError, fall.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (r *FeedbackRepository) GetModelFeedback(ctx context.Context, modelId int, order string) (*model.ModelFeedbackResponse, fall.Error) {
	query := fmt.Sprintf(`
	SELECT f.feedback_id as f_id, f.created_at as created_at, f.updated_at as updated_at, f.feedback_text as f_text,
	f.rate as f_rate,
	(select avg(f.rate) as avg_rate from feedback as f 
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	WHERE pm.product_model_id = $1 AND f.is_hidden = FALSE
	) as avg_rate,
	(select count(f.feedback_id) from feedback as f 
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	WHERE pm.product_model_id = $1 AND f.is_hidden = FALSE
	) as rate_count,
	f.product_model_id as f_model_id, f.is_hidden as f_hidden,
	u.user_id as u_id, u.email as u_email,
	u.avatar_path as u_avatar_path, u.first_name as u_first_name, u.last_name as u_last_name
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	WHERE pm.product_model_id = $1 AND f.is_hidden = FALSE
	ORDER BY f.updated_at %s;`, order)
	rows, err := r.db.Query(ctx, query, modelId)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()
	var feedbackResponse model.ModelFeedbackResponse
	var feedback []model.Feedback
	var founded bool = false
	for rows.Next() {

		f := model.Feedback{}

		err := rows.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &feedbackResponse.AvgRate, &feedbackResponse.RateCount, &f.ModelId, &f.Hidden,
			&f.User.Id, &f.User.Email, &f.User.Avatar, &f.User.FirstName, &f.User.LastName,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		feedback = append(feedback, f)
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.FeedbackNotFound, fall.STATUS_NOT_FOUND)
	}

	feedbackResponse.Feedback = feedback

	return &feedbackResponse, nil
}

func (r *FeedbackRepository) GetAll(ctx context.Context, order string) ([]model.Feedback, fall.Error) {
	query := fmt.Sprintf(`
	SELECT f.feedback_id as f_id, f.created_at as created_at, f.updated_at as updated_at, f.feedback_text as f_text, f.rate as f_rate,
	f.product_model_id as f_model_id, f.is_hidden as f_hidden,
	u.user_id as u_id, u.email as u_email,
	u.avatar_path as u_avatar_path, u.first_name as u_first_name, u.last_name as u_last_name
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	ORDER BY f.updated_at %s;
	`, order)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()
	var result []model.Feedback
	for rows.Next() {

		f := model.Feedback{}

		err := rows.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &f.ModelId, &f.Hidden,
			&f.User.Id, &f.User.Email, &f.User.Avatar, &f.User.FirstName, &f.User.LastName,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		result = append(result, f)
	}

	if rows.Err() != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	return result, nil
}

func (r *FeedbackRepository) DeleteFeedback(ctx context.Context, feedbackId int) fall.Error {
	query := `
	DELETE FROM feedback WHERE feedback_id = $1;
	`
	_, err := r.db.Exec(ctx, query, feedbackId)

	if err != nil {
		return fall.NewErr(msg.FeedbackCreateError, fall.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (r *FeedbackRepository) ToggleHidden(ctx context.Context, feedbackId int) fall.Error {
	query := `
	SELECT is_hidden FROM feedback WHERE feedback_id = $1;
	`
	var isHidden bool

	row := r.db.QueryRow(ctx, query, feedbackId)

	err := row.Scan(&isHidden)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	updateQuery := `UPDATE feedback SET is_hidden = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE feedback_id = $2;`

	_, err = r.db.Exec(ctx, updateQuery, !isHidden, feedbackId)
	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil

}
