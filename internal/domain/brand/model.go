package brand

type Brand struct {
	Id          int     `json:"id" example:"2"`
	Title       string  `json:"title" example:"adidas"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	ImgPath     *string `json:"img_path"`
}

type CreateBrandDto struct {
	Title       string  `json:"title" example:"adidas" validate:"required,min=2"`
	Description *string `json:"description"`
	ImgPath     *string `json:"img_path"`
}
