package model

type CreateOptionDto struct {
	Title string `json:"title" validate:"required,min=3" example:"Цвет"`
	Slug  string `json:"slug" validate:"required,min=3" example:"color"`
}

type UpdateOptionDto struct {
	Title      *string `json:"title" validate:"omitempty,min=3"`
	ForCatalog *bool   `json:"for_catalog" validate:"omitempty"`
	Slug       *string `json:"slug" validate:"omitempty,min=3"`
}

type UpdateOptionValueDto struct {
	Value *string `json:"value" example:"Желтый" validate:"omitempty,min=1"`
	Info  *string `json:"info" validate:"omitempty,min=1"`
}

type CreateOptionValueDto struct {
	Value    string  `json:"value" example:"Желтый" validate:"required,min=1"`
	Info     *string `json:"info"`
	OptionId int     `json:"option_id" example:"4" validate:"required,min=1"`
}

type AddOptionToProductModelDto struct {
	ProductModelId int `json:"product_model_id" validate:"required,min=1"`
	OptionId       int `json:"option_id" validate:"required,min=1"`
	ValueId        int `json:"value_id" validate:"required,min=1"`
}

type CreateSizeDto struct {
	Value string `json:"value" validate:"required,min=1,max=10" example:"42"`
}

type AddSizeToProductModelDto struct {
	ProductModelId int    `json:"product_model_id" validate:"required,min=1"`
	SizeId         int    `json:"size_id" validate:"required,min=1"`
	InStock        int    `json:"in_stock" validate:"required,min=0"`
	Literal        string `json:"literal" validate:"required"`
}

type Option struct {
	Id         int           `json:"id" example:"4" validate:"required"`
	Title      string        `json:"title" example:"Цвет" validate:"required"`
	Slug       string        `json:"slug" example:"color" validate:"required"`
	ForCatalog bool          `json:"for_catalog" validate:"required"`
	Values     []OptionValue `json:"values"`
}

type OptionValue struct {
	Id       *int    `json:"id" example:"44"`
	Value    *string `json:"value" example:"Желтый"`
	Info     *string `json:"info"`
	OptionId *int    `json:"option_id" example:"4"`
}

type Size struct {
	Id      int    `json:"id" example:"44" validate:"required"`
	Numeric string `json:"numeric" validate:"required"`
	Literal string `json:"literal" validate:"required"`
}

type CatalogOption struct {
	Id     int            `json:"option_id" validate:"required"`
	Title  string         `json:"title" validate:"required"`
	Slug   string         `json:"slug" validate:"required"`
	Values []CatalogValue `json:"values" validate:"required"`
}

type CatalogValue struct {
	Id       int    `json:"value_id" validate:"required"`
	Value    string `json:"value" validate:"required"`
	OptionId int    `json:"option_id" validate:"required"`
}

type CatalogSize struct {
	Id    int    `json:"size_id" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type CatalogBrand struct {
	Id    int    `json:"brand_id" validate:"required"`
	Title string `json:"brand_title" validate:"required"`
}

type CatalogPrice struct {
	Max int `json:"max_price" validate:"required"`
	Min int `json:"min_price" validate:"required"`
}

type CatalogFilters struct {
	Options []CatalogOption `json:"options" validate:"required"`
	Sizes   []CatalogSize   `json:"sizes" validate:"required"`
	Brands  []CatalogBrand  `json:"brands" validate:"required"`
	Price   CatalogPrice    `json:"price" validate:"required"`
}
