package model

type Brand struct {
	Id          int     `json:"id" example:"2" validate:"required"`
	Title       string  `json:"title" example:"adidas" validate:"required"`
	Slug        string  `json:"slug" example:"adidas" validate:"required"`
	Description *string `json:"description" example:"бренд одежды..."`
	ImgPath     *string `json:"img_path" example:"/static/category/adidas.webp"`
}

type CreateBrandDto struct {
	Title       string  `json:"title" example:"adidas" validate:"required,min=2"`
	Description *string `json:"description" validate:"omitempty,min=10" example:"бренд одежды"`
	ImgPath     *string `json:"img_path" validate:"omitempty,filepath" example:"/static/category/adidas.webp"`
}

type UpdateBrandDto struct {
	Title       *string `json:"title" example:"adidas" validate:"omitempty,min=2" `
	Description *string `json:"description" validate:"omitempty,min=10" example:"бренд одежды"`
	ImgPath     *string `json:"img_path" validate:"omitempty,filepath" example:"/static/category/adidas.webp"`
}
