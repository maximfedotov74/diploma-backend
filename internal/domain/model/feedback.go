package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type FeedbakAdminFilterEnum string

const (
	All        FeedbakAdminFilterEnum = "all"
	OnlyActive FeedbakAdminFilterEnum = "only_active"
	OnlyHidden FeedbakAdminFilterEnum = "only_hidden"
)

func FeedbakAdminFilterEnumValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(All), string(OnlyActive), string(OnlyHidden):
		return true
	}
	return false
}

type ModelFeedbackResponse struct {
	Feedback  []Feedback `json:"feedback" validate:"required"`
	AvgRate   *float32   `json:"avg_rate"`
	RateCount *int       `json:"rate_count"`
}

type Feedback struct {
	Id        int          `json:"id" example:"2" validate:"required"`
	CreatedAt time.Time    `json:"created_at" validate:"required"`
	UpdatedAt time.Time    `json:"updated_at" validate:"required"`
	Text      string       `json:"text" validate:"required,min=3" example:"Хороший товар"`
	Rate      int8         `json:"rate" validate:"required,min=1,max=5" example:"3"`
	ModelId   int          `json:"model_id" validate:"required,min=1" example:"4"`
	Hidden    bool         `json:"is_hidden" validate:"required"`
	User      FeedbackUser `json:"user" validate:"required"`
}

type FeedbackUser struct {
	Id        int     `json:"id" validate:"required"`
	Email     string  `json:"email" validate:"required"`
	Avatar    *string `json:"avatar"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type AddFeedbackDto struct {
	Text    string `json:"text" validate:"required,min=3" example:"Хороший товар"`
	Rate    int8   `json:"rate" validate:"required,min=1,max=5" example:"3"`
	ModelId int    `json:"model_id" validate:"required,min=1" example:"4"`
}

type AdminAllFeedbackResponse struct {
	Feedback []Feedback `json:"feedback" validate:"required"`
	Total    int        `json:"total" validate:"required"`
}
