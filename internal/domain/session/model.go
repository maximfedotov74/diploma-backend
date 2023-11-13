package session

type Session struct {
	SessionID int    `json:"session_id" db:"session_id"`
	UserId    int    `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	Token     string `json:"token" db:"token"`
}
