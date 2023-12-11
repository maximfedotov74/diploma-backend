package category

type CreateCategoryDto struct {
	Title      string  `json:"title" db:"title" example:"Мужская Верхняя одежда" validate:"required,min=3"`
	ShortTitle string  `json:"short_title" db:"short_title" example:"Верхняя одежда" validate:"required,min=3"`
	ImgPath    *string `json:"img_path" db:"img_path" example:"/static/example.webp" validate:"omitempty,filepath"`
	ParentId   *int    `json:"parent_category_id" validate:"omitempty" example:"4"`
}

type UpdateCategoryDto struct {
	ImgPath    *string `json:"img_path" db:"img_path" example:"/static/example.webp" validate:"omitempty,filepath"`
	Title      *string `json:"title" db:"title" example:"Мужская Верхняя одежда" validate:"omitempty,min=3"`
	ShortTitle *string `json:"short_title" db:"short_title" example:"Верхняя одежда" validate:"omitempty,min=3"`
}
