package model

type UserSettings struct {
	Id                    int     `json:"user_settings_id" db:"user_settings_id" example:"1"`
	ActivationAccountLink *string `json:"activation_account_link" db:"activation_account_link"`
	AuthProvider          string  `json:"auth_provider" db:"auth_provider"`
	UserId                int     `json:"user_id" db:"user_id"`
}
