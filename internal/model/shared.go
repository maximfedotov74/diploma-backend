package model

type СompletedOperation struct {
	Completed bool `json:"completed" example:"true"`
}

type UserContextData struct {
	User      User   `json:"user"`
	UserAgent string `json:"user_agent"`
}
