package msg

const (
	ProductCreateError                 = "Ошибка при создании товара!"
	ProductCreateModelError            = "Ошибка при создании модели товара!"
	ProductNotFound                    = "Товар не найден!"
	ProductModelNotFound               = "Модель товара не найдена!"
	ProductAddPhotoError               = "Ошибка при добавлении фотографии!"
	ProductUpdateError                 = "Ошибка при обновлении товара!"
	ProductModelUpdateError            = "Ошибка при обновлении модели товара!"
	ProductModelSlugUnique             = "Slug модели должен быть уникален!"
	ProductInStockCannotBeLessThanZero = "Количество товара на складе не может быть меньше 0"
	ProductDeleteError                 = "Ошибка при удалении товара!"
	ProductModelDeleteError            = "Ошибка при удалении модели товара!"
)
