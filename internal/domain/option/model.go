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

type Size struct {
	Id      int    `json:"id" example:"44"`
	Numeric string `json:"numeric"`
	Literal string `json:"literal"`
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

type CatalogOption struct {
	Id     int            `json:"option_id"`
	Title  string         `json:"title"`
	Slug   string         `json:"slug"`
	Values []CatalogValue `json:"values"`
}

type CatalogValue struct {
	Id       int    `json:"value_id"`
	Value    string `json:"value"`
	OptionId int    `json:"option_id"`
}

type CatalogSize struct {
	Id    int    `json:"size_id"`
	Value string `json:"value"`
}

type CatalogFilters struct {
	Options []CatalogOption `json:"options"`
	Sizes   []CatalogSize   `json:"sizes"`
}

type CatalogSizeSorter []CatalogSize

func (p CatalogSizeSorter) Len() int           { return len(p) }
func (p CatalogSizeSorter) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p CatalogSizeSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
