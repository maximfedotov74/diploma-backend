package auth

type LoginDto struct {
	Email    string `json:"email" validate:"required,email" example:"makc@mail.ru"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}
