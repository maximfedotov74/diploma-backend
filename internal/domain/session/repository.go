package session

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sr *SessionRepository) FindByAgentAndToken(agent string, token string) (*Session, exception.Error) {

	ctx := context.Background()

	query := "SELECT session_id, user_id, user_agent, token FROM public.session WHERE user_agent = $1 AND token = $2;"

	sessionModel := Session{}

	row := sr.db.QueryRow(ctx, query, agent, token)

	err := row.Scan(&sessionModel.SessionID, &sessionModel.UserId, &sessionModel.UserAgent, &sessionModel.Token)

	if err != nil {
		return nil, exception.NewErr(sessionNotFound, exception.STATUS_NOT_FOUND)
	}

	return &sessionModel, nil

}

func (tr *SessionRepository) RemoveSession(token string, agent string) exception.Error {
	query := "DELETE FROM session WHERE token = $1 AND user_agent = $2;"
	_, err := tr.db.Exec(context.Background(), query, token, agent)
	if err != nil {
		return exception.ServerError(err.Error())
	}
	return nil
}

func (tr *SessionRepository) RemoveExceptCurrentSession(userId int, agent string) exception.Error {
	query := "DELETE FROM session WHERE user_id = $1 AND user_agent != $2;"
	_, err := tr.db.Exec(context.Background(), query, userId, agent)
	if err != nil {
		return exception.ServerError(err.Error())
	}
	return nil
}

func (tr *SessionRepository) CreateSession(dto CreateSessionDto) exception.Error {

	ctx := context.Background()

	query := "SELECT session_id, user_id, user_agent, token FROM public.session WHERE user_id = $1 AND user_agent = $2;"

	sessionModel := Session{}

	row := tr.db.QueryRow(ctx, query, dto.UserId, dto.UserAgent)

	err := row.Scan(&sessionModel.SessionID, &sessionModel.UserId, &sessionModel.UserAgent, &sessionModel.Token)

	if err != nil {
		if err == pgx.ErrNoRows {
			query = "INSERT INTO session (token, user_agent, user_id) VALUES ($1, $2, $3) RETURNING session_id;"
			_, err = tr.db.Exec(ctx, query, dto.Token, dto.UserAgent, dto.UserId)
			if err != nil {
				return exception.NewErr(sessionCreateError, exception.STATUS_INTERNAL_ERROR)
			}
			return nil
		}
		return exception.ServerError(err.Error())
	}
	query = "UPDATE public.session SET token = $1, updated_at = CURRENT_TIMESTAMP WHERE session_id = $2;"

	_, err = tr.db.Exec(ctx, query, dto.Token, sessionModel.SessionID)

	if err != nil {

		return exception.NewErr(sessionUpdateError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil

}
