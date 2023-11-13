package user

type ChangePasswordDto struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
	Code        string `json:"code" validate:"required,min=6,max=6" example:"123456"`
}

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email" example:"makc@mail.ru"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}
