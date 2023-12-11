package feedback

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type FeedbackRepository struct {
	db *pgxpool.Pool
}

func NewFeedbackRepository(db *pgxpool.Pool) *FeedbackRepository {
	return &FeedbackRepository{
		db: db,
	}
}

func (fr *FeedbackRepository) AddFeedback(userId int, dto AddFeedbackDto) exception.Error {
	query := `
	INSERT INTO feedback (feedback_text, rate, product_model_id, user_id) VALUES ($1,$2,$3,$4);
	`
	_, err := fr.db.Exec(context.Background(), query, dto.Text, dto.Rate, dto.ModelId, userId)

	if err != nil {
		return exception.NewErr(feedbackCreateError, exception.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (fr *FeedbackRepository) GetModelFeedback(modelId int, order string) (*ModelFeedbackResponse, exception.Error) {
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
	u.avatar_path as u_avatar_path, us.first_name as u_first_name, us.last_name as u_last_name, 
	us.patronymic as u_patronymic
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN user_settings as us ON us.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	WHERE pm.product_model_id = $1 AND f.is_hidden = FALSE
	ORDER BY f.updated_at %s;`, order)
	rows, err := fr.db.Query(context.Background(), query, modelId)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()
	var feedbackResponse ModelFeedbackResponse
	var feedback []Feedback
	var founded bool = false
	for rows.Next() {

		f := Feedback{}

		err := rows.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &feedbackResponse.AvgRate, &feedbackResponse.RateCount, &f.ModelId, &f.Hidden,
			&f.User.Id, &f.User.Email, &f.User.Avatar, &f.User.FirstName, &f.User.LastName,
			&f.User.Patronymic,
		)
		if err != nil {
			return nil, exception.ServerError(err.Error())
		}
		feedback = append(feedback, f)
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, exception.NewErr(feedbackNotFound, exception.STATUS_NOT_FOUND)
	}

	feedbackResponse.Feedback = feedback

	return &feedbackResponse, nil
}

func (fr *FeedbackRepository) GetAll(order string) ([]Feedback, exception.Error) {
	query := fmt.Sprintf(`
	SELECT f.feedback_id as f_id, f.created_at as created_at, f.updated_at as updated_at, f.feedback_text as f_text, f.rate as f_rate,
	f.product_model_id as f_model_id, f.is_hidden as f_hidden,
	u.user_id as u_id, u.email as u_email,
	u.avatar_path as u_avatar_path, us.first_name as u_first_name, us.last_name as u_last_name, 
	us.patronymic as u_patronymic
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN user_settings as us ON us.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	ORDER BY f.updated_at %s;
	`, order)
	rows, err := fr.db.Query(context.Background(), query)
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()
	var result []Feedback
	for rows.Next() {

		f := Feedback{}

		err := rows.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &f.ModelId, &f.Hidden,
			&f.User.Id, &f.User.Email, &f.User.Avatar, &f.User.FirstName, &f.User.LastName,
			&f.User.Patronymic,
		)
		if err != nil {
			return nil, exception.ServerError(err.Error())
		}
		result = append(result, f)
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	return result, nil
}

func (fr *FeedbackRepository) DeleteFeedback(feedbackId int) exception.Error {
	query := `
	DELETE FROM feedback WHERE feedback_id = $1;
	`
	_, err := fr.db.Exec(context.Background(), query, feedbackId)

	if err != nil {
		return exception.NewErr(feedbackCreateError, exception.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (fr *FeedbackRepository) ToggleHidden(feedbackId int) exception.Error {
	query := `
	SELECT is_hidden FROM feedback WHERE feedback_id = $1;
	`
	var isHidden bool

	row := fr.db.QueryRow(context.Background(), query, feedbackId)

	err := row.Scan(&isHidden)

	if err != nil {
		return exception.ServerError(err.Error())
	}

	updateQuery := `UPDATE feedback SET is_hidden = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE feedback_id = $2;`

	_, err = fr.db.Exec(context.Background(), updateQuery, !isHidden, feedbackId)
	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil

}
