package messages

const (
	USER_NOT_FOUND                 = "Пользователь не найден!"
	USER_EXISTS                    = "Пользователь с таким email уже существует!"
	INVALID_CREDENTIALS            = "Неверный логин или пароль!"
	ACTIVATION_ERROR               = "Произошла ошибка при активации! Попробуйте в другой раз."
	ACTIVATION_NOT_FOUND           = "Пользователя с такой ссылкой не существует!"
	UNAUTHORIZED                   = "Неавторизован!"
	FORBIDDEN                      = "Доступ запрещен!"
	UPDATE_PASSWORD_ERROR          = "Ошибка при смене пароля!"
	BAD_PASSWORD                   = "Введенный пароль не совпадает с текущим!"
	BAD_NEW_PASSWORD               = "Новый пароль должен отличаться от старого!"
	CHANGE_PASSWORD_CODE_NOT_FOUND = "Неверный код или вышел его срок действия!"
)
