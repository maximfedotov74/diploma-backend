package feedback

import "time"

type ModelFeedbackResponse struct {
	Feedback  []Feedback `json:"feedback"`
	AvgRate   *float32   `json:"avg_rate"`
	RateCount *int       `json:"rate_count"`
}

type Feedback struct {
	Id        int          `json:"id" example:"2"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Text      string       `json:"text" validate:"required,min=3" example:"Хороший товар"`
	Rate      int8         `json:"rate" validate:"required,min=1,max=5" example:"3"`
	ModelId   int          `json:"model_id" validate:"required,min=1" example:"4"`
	Hidden    bool         `json:"is_hidden"`
	User      FeedbackUser `json:"user"`
}

type FeedbackUser struct {
	Id         int     `json:"id"`
	Email      string  `json:"email"`
	Avatar     *string `json:"avatar"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Patronymic *string `json:"patronymic"`
}
