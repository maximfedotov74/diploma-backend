package user

const (
	userNotFound               = "Пользователь не найден!"
	userExists                 = "Пользователь с таким email уже существует!"
	activationError            = "Произошла ошибка при активации! Попробуйте в другой раз."
	activationNotFound         = "Пользователя с такой ссылкой не существует!"
	updatePasswordError        = "Ошибка при смене пароля!"
	badPassword                = "Введенный пароль не совпадает с текущим!"
	badNewPassword             = "Новый пароль должен отличаться от старого!"
	changePasswordCodeNotFound = "Неверный код или вышел его срок действия!"
	changePasswrodError        = "Ошибка при смене пароля!"
	createChangeCodeError      = "Ошибка при создании кода для смены пароля!"
)
