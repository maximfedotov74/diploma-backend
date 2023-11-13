package session

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sr *SessionRepository) FindByAgentAndToken(agent string, token string) (*Session, error) {

	ctx := context.Background()

	query := "SELECT session_id, user_id, user_agent, token FROM public.session WHERE user_agent = $1 AND token = $2;"

	sessionModel := Session{}

	row := sr.db.QueryRow(ctx, query, agent, token)

	err := row.Scan(&sessionModel.SessionID, &sessionModel.UserId, &sessionModel.UserAgent, &sessionModel.Token)

	if err != nil {
		return nil, err
	}

	return &sessionModel, nil

}

func (tr *SessionRepository) RemoveSession(token string) error {
	query := "DELETE FROM session WHERE token = $1;"
	_, err := tr.db.Exec(context.Background(), query, token)
	if err != nil {
		return err
	}
	return nil
}

func (tr *SessionRepository) CreateSession(dto CreateSessionDto) error {

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
				return err
			}

			return nil
		} else {

			return err
		}
	}
	query = "UPDATE public.session SET token = $1, updated_at = CURRENT_TIMESTAMP WHERE session_id = $2;"

	_, err = tr.db.Exec(ctx, query, dto.Token, sessionModel.SessionID)

	if err != nil {

		return err
	}

	return nil

}
