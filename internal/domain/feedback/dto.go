package feedback

type AddFeedbackDto struct {
	Text    string `json:"text" validate:"required,min=3" example:"Хороший товар"`
	Rate    int8   `json:"rate" validate:"required,min=1,max=5" example:"3"`
	ModelId int    `json:"model_id" validate:"required,min=1" example:"4"`
}
