package model

type ChangePasswordCode struct {
	ChangePasswordCodeId int    `json:"change_password_code_id" db:"change_password_code_id"`
	Code                 string `json:"code" db:"code"`
	UserId               int    `json:"user_id" db:"user_id"`
}
