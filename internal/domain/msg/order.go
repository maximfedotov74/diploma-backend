package msg

const (
	OrderDeliveryPointConditionConflict = "Условия заказа не совпадают с регламентом пункта выдачи!"
	OrderErrorWhenAddModelsToProduct    = "Произошла ошибка при добавлени товаров к заказу!"
	OrderErrorWhenAddActivationLink     = "Ошибка при создании ссылки активации!"
	OrderErrorWhenAddDeliveryPoint      = "Ошибка при добавлении пункта выдачи!"
	OrderAlreadyActivated               = "Заказ уже активирован!"
	OrderErrorWhenActivate              = "Ошибка при активации заказа!"
	OrderNotFound                       = "Заказ не найден!"
	OrderActivationLinkNotFound         = "Ссылка активации заказа не найдена!"
	OrderErrorWhenCancel                = "Ошибка при отмене заказа!"
	OrderErrorWhenChangeStatus          = "Ошибка при смене статуса!"
	OrderErrorWhenChangeDeliveryDate    = "Ошибка при смене даты доставки!"
)
