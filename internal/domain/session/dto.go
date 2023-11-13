package session

type CreateSessionDto struct {
	UserId    int    `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	Token     string `json:"token" db:"token"`
}
