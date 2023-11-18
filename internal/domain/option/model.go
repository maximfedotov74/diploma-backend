package option

type Option struct {
	Id     int           `json:"id" example:"4"`
	Title  string        `json:"title" example:"Цвет"`
	Slug   string        `json:"slug" example:"color"`
	Values []OptionValue `json:"values"`
}

type OptionValue struct {
	Id       *int    `json:"id" example:"44"`
	Value    *string `json:"value" example:"Желтый"`
	Info     *string `json:"info"`
	OptionId *int    `json:"option_id" example:"4"`
}

type ProductModelOption struct {
	Id             int
	ProductModelId int
	OptionId       int
	ValueId        int
}

type CategoryOption struct {
	Id         int
	CategoryId int
	OptionId   int
}
