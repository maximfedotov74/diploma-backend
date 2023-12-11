package brand

type CreateBrandDto struct {
	Title       string  `json:"title" example:"adidas" validate:"required,min=2" example:"адидас"`
	Description *string `json:"description" validate:"omitempty,min=10" example:"бренд одежды"`
	ImgPath     *string `json:"img_path" validate:"omitempty,filepath" example:"/static/category/adidas.webp"`
}

type UpdateBrandDto struct {
	Title       *string `json:"title" example:"adidas" validate:"omitempty,min=2" example:"адидас"`
	Description *string `json:"description" validate:"omitempty,min=10" example:"бренд одежды"`
	ImgPath     *string `json:"img_path" validate:"omitempty,filepath" example:"/static/category/adidas.webp"`
}
