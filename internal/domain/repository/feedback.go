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

func (r *FeedbackRepository) GetAll(ctx context.Context, order string, page int, filter string) (*model.AdminAllFeedbackResponse, fall.Error) {

	limit := 16

	offset := page*limit - limit
	pagination := fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	where := ""

	switch filter {
	case string(model.All):
		where = ""
	case string(model.OnlyActive):
		where = "WHERE f.is_hidden = FALSE"
	case string(model.OnlyHidden):
		where = "WHERE f.is_hidden = TRUE"
	default:
		where = ""
	}

	query := fmt.Sprintf(`
	SELECT f.feedback_id as f_id, f.created_at as created_at, f.updated_at as updated_at, f.feedback_text as f_text, f.rate as f_rate,
	f.product_model_id as f_model_id, f.is_hidden as f_hidden,
	u.user_id as u_id, u.email as u_email,
	u.avatar_path as u_avatar_path, u.first_name as u_first_name, u.last_name as u_last_name,
	(select count(distinct f.feedback_id) 
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	%s
	) as total_count
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	%s
	ORDER BY f.created_at %s %s;
	`, where, where, order, pagination)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()
	var result []model.Feedback
	var total int
	for rows.Next() {

		f := model.Feedback{}

		err := rows.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &f.ModelId, &f.Hidden,
			&f.User.Id, &f.User.Email, &f.User.Avatar, &f.User.FirstName, &f.User.LastName, &total,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		result = append(result, f)
	}

	if rows.Err() != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	return &model.AdminAllFeedbackResponse{
		Feedback: result,
		Total:    total,
	}, nil
}

func (r *FeedbackRepository) GetMyFeedback(ctx context.Context, userId int) ([]model.UserFeedback, fall.Error) {
	q := `
 	SELECT f.feedback_id as f_id, f.created_at as created_at, f.updated_at as updated_at, f.feedback_text as f_text, f.rate as f_rate,
	f.is_hidden as f_hidden, p.product_id as p_id, p.title as p_title, p.description as p_descr,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path, b.description as d_description,
	ct.category_id as ct_id, ct.title as ct_title,
	ct.short_title as ct_short_title, ct.slug as ct_slug,
	ct.img_path as c_img_path, ct.parent_category_id as c_parent_id,
	pm.product_model_id as model_id, pm.product_id as m_product_id, pm.slug as m_slug, pm.article as m_article, pm.price as model_price, pm.discount as model_discount,
	pm.main_image_path as pm_main_img
	FROM feedback as f
	INNER JOIN public.user as u ON f.user_id = u.user_id
	INNER JOIN product_model as pm ON pm.product_model_id = f.product_model_id
	INNER JOIN product as p ON pm.product_id = p.product_id
	INNER JOIN category ct ON p.category_id = ct.category_id 
	INNER JOIN brand b on p.brand_id = b.brand_id
	WHERE u.user_id = $1
	ORDER BY f.created_at DESC;
 `
	rows, err := r.db.Query(ctx, q, userId)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()
	var result []model.UserFeedback
	for rows.Next() {

		f := model.UserFeedback{}
		p := model.Product{}
		m := model.ProductModel{}

		err := rows.Scan(&f.Id, &f.CreatedAt, &f.UpdatedAt, &f.Text, &f.Rate, &f.Hidden, &p.Id, &p.Title, &p.Description,
			&p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath, &p.Brand.Description,
			&p.Category.Id, &p.Category.Title, &p.Category.ShortTitle, &p.Category.Slug, &p.Category.ImgPath, &p.Category.ParentId,
			&m.Id, &m.ProductId, &m.Slug, &m.Article, &m.Price, &m.Discount,
			&m.ImagePath)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		f.Product = p
		f.Model = m
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
