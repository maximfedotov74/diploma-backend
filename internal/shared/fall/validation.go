package fall

import (
	"github.com/go-playground/validator/v10"
)

type validationErrorItem struct {
	Key     string `json:"key" example:"email"`
	Message string `json:"message" example:"email is invalid"`
}

type ValidationError struct {
	Status  int                   `json:"status" example:"400" validate:"required"`
	Errors  []validationErrorItem `json:"errors"`
	Message string                `json:"message" validate:"required"`
}

func NewValidErr(e []validationErrorItem) ValidationError {
	return ValidationError{Status: 400, Errors: e, Message: "Ошибка при валидации данных!"}
}

func error_message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return VALIDATION_REQUIRED
	case "email":
		return VALIDATION_EMAIL
	case "min":
		return VALIDATION_MIN
	case "max":
		return VALIDATION_MAX
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
