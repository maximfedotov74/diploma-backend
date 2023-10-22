package lib

import "github.com/go-playground/validator/v10"

type validationErrorItem struct {
	Key     string `json:"key" example:"email"`
	Message string `json:"message" example:"email is invalid"`
}

type ValidationError struct {
	Status int                   `json:"status" example:"400"`
	Errors []validationErrorItem `json:"errors"`
}

func NewValidErr(e []validationErrorItem) ValidationError {
	return ValidationError{Status: 400, Errors: e}
}

func error_message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Это обязательное поле"
	case "email":
		return "Некорректный адрес электронной почты!"
	case "min":
		return "Значение меньше минимального предела"
	case "max":
		return "Значение больше максимального предела"
	}

	return fe.Error()
}

func ValidationMessages(ve validator.ValidationErrors) []validationErrorItem {
	out := make([]validationErrorItem, len(ve))

	for idx, err := range ve {
		out[idx] = validationErrorItem{err.Field(), error_message(err)}
	}
	return out
}
