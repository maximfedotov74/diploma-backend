package option

type CreateOptionDto struct {
	Title string `json:"title" validate:"required,min=3" example:"Цвет"`
	Slug  string `json:"slug" validate:"required,min=3" example:"color"`
}

type UpdateOptionDto struct {
	Title string `json:"title" validate:"required,min=3" example:"Цвет"`
}

type CreateOptionValueDto struct {
	Value    string  `json:"value" example:"Желтый"`
	Info     *string `json:"info"`
	OptionId int     `json:"option_id" example:"4"`
}

type AddOptionToProductModelDto struct {
	ProductModelId int `json:"product_model_id" validate:"required,min=1"`
	OptionId       int `json:"option_id" validate:"required,min=1"`
	ValueId        int `json:"value_id" validate:"required,min=1"`
}
