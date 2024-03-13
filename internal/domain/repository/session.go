package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type SessionRepository struct {
	db db.PostgresClient
}

func NewSessionRepository(db db.PostgresClient) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, dto model.CreateSessionDto) fall.Error {
	q := "SELECT session_id, user_id, user_agent, token FROM public.session WHERE user_id = $1 AND user_agent = $2;"

	sessionModel := model.Session{}

	row := r.db.QueryRow(ctx, q, dto.UserId, dto.UserAgent)

	err := row.Scan(&sessionModel.SessionId, &sessionModel.UserId, &sessionModel.UserAgent, &sessionModel.Token)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			q = "INSERT INTO session (token, user_agent, user_id) VALUES ($1, $2, $3) RETURNING session_id;"
			_, err = r.db.Exec(ctx, q, dto.Token, dto.UserAgent, dto.UserId)
			if err != nil {
				return fall.NewErr(msg.SessionCreateError, fall.STATUS_INTERNAL_ERROR)
			}
			return nil
		}
		return fall.ServerError(err.Error())
	}
	q = "UPDATE public.session SET token = $1, updated_at = CURRENT_TIMESTAMP WHERE session_id = $2;"

	_, err = r.db.Exec(ctx, q, dto.Token, sessionModel.SessionId)

	if err != nil {

		return fall.NewErr(msg.SessionUpdateError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (r *SessionRepository) FindByAgentAndToken(ctx context.Context, agent string, token string) (*model.Session, fall.Error) {

	query := "SELECT session_id, user_id, user_agent, token, created_at, updated_at FROM public.session WHERE user_agent = $1 AND token = $2;"

	sessionModel := model.Session{}

	row := r.db.QueryRow(ctx, query, agent, token)

	err := row.Scan(&sessionModel.SessionId, &sessionModel.UserId, &sessionModel.UserAgent,
		&sessionModel.Token, &sessionModel.CreatedAt, &sessionModel.UpdatedAt,
	)

	if err != nil {
		return nil, fall.NewErr(msg.SessionNotFound, fall.STATUS_NOT_FOUND)
	}

	return &sessionModel, nil

}

func (r *SessionRepository) FindByAgentAndUserId(ctx context.Context, agent string, userId int) (*model.Session, fall.Error) {

	query := "SELECT session_id, user_id, user_agent, token, created_at, updated_at FROM public.session WHERE user_agent = $1 AND user_id = $2;"

	sessionModel := model.Session{}

	row := r.db.QueryRow(ctx, query, agent, userId)

	err := row.Scan(&sessionModel.SessionId, &sessionModel.UserId, &sessionModel.UserAgent,
		&sessionModel.Token, &sessionModel.CreatedAt, &sessionModel.UpdatedAt,
	)

	if err != nil {
		return nil, fall.NewErr(msg.SessionNotFound, fall.STATUS_NOT_FOUND)
	}

	return &sessionModel, nil

}

func (r *SessionRepository) RemoveSession(ctx context.Context, userId int, sessionId int) fall.Error {
	query := "DELETE FROM session WHERE user_id = $1 AND session_id = $2;"
	_, err := r.db.Exec(ctx, query, userId, sessionId)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) fall.Error {
	query := "DELETE FROM session WHERE user_id = $1 AND session_id != $2;"
	_, err := r.db.Exec(ctx, query, userId, sessionId)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) RemoveSessionByToken(ctx context.Context, token string) fall.Error {
	query := "DELETE FROM session WHERE token = $1;"
	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) RemoveAllSessions(ctx context.Context, userId int) fall.Error {
	query := "DELETE FROM session WHERE user_id = $1;"

	_, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) GetUserSessions(ctx context.Context, userId int, token string) (*model.UserSessionsResponse, fall.Error) {
	q := `SELECT session_id, user_id, user_agent, token, created_at, updated_at
	FROM public.session WHERE user_id = $1 ORDER BY updated_at DESC;`

	var sessions []model.Session

	rows, err := r.db.Query(ctx, q, userId)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	response := model.UserSessionsResponse{}

	for rows.Next() {
		s := model.Session{}

		err := rows.Scan(&s.SessionId, &s.UserId, &s.UserAgent,
			&s.Token, &s.CreatedAt, &s.UpdatedAt,
		)

		if err != nil {
			return nil, fall.NewErr(msg.SessionNotFound, fall.STATUS_NOT_FOUND)
		}

		if s.Token == token {
			response.Current = &s
		}

		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	response.All = sessions

	return &response, nil
}
