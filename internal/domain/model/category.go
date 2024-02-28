package model

type CategoryModel struct {
	Id         int     `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title      string  `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug       string  `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle string  `json:"short_title" validate:"required"`
	ImgPath    *string `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId   *int    `json:"parent_category_id" db:"parent_category_id" example:"2"`
}

type Child struct {
	Id            int                 `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title         string              `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug          string              `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle    string              `json:"short_title" validate:"required"`
	ImgPath       *string             `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int                `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8               `json:"level" validate:"required"`
	Subcategories []*CategoryRelation `json:"subcategories" validate:"required"`
}

type CategoryRelation struct {
	Id            int      `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title         string   `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug          string   `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle    string   `json:"short_title" validate:"required"`
	ImgPath       *string  `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int     `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8    `json:"level" validate:"required"`
	Subcategories []*Child `json:"subcategories" validate:"required"`
}

type Category struct {
	Id            int         `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title         string      `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug          string      `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle    string      `json:"short_title" validate:"required"`
	ImgPath       *string     `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int        `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8       `json:"level" validate:"required"`
	Subcategories []*Category `json:"subcategories" validate:"required"`
}

type CatalogChild struct {
	Id            int                        `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title         string                     `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug          string                     `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle    string                     `json:"short_title" validate:"required"`
	ImgPath       *string                    `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int                       `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8                      `json:"level" validate:"required"`
	Active        bool                       `json:"active" validate:"required"`
	Subcategories []*CatalogCategoryRelation `json:"subcategories" validate:"required"`
}

type CatalogCategoryResponse struct {
	CatalogCategories СatalogCategory `json:"catalog_categories" validate:"required"`
	Current           CategoryModel   `json:"current" validate:"required"`
}
type CatalogCategoryRelationResponse struct {
	CatalogCategories CatalogCategoryRelation `json:"catalog_categories" validate:"required"`
	Current           CategoryModel           `json:"current" validate:"required"`
}

type CatalogCategoryRelation struct {
	Id            int             `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title         string          `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug          string          `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle    string          `json:"short_title" validate:"required"`
	ImgPath       *string         `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int            `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8           `json:"level" validate:"required"`
	Active        bool            `json:"active" validate:"required"`
	Subcategories []*CatalogChild `json:"subcategories" validate:"required"`
}

type СatalogCategory struct {
	Id            int                `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title         string             `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	Slug          string             `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
	ShortTitle    string             `json:"short_title" validate:"required"`
	ImgPath       *string            `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int               `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8              `json:"level" validate:"required"`
	Active        bool               `json:"active" validate:"required"`
	Subcategories []*СatalogCategory `json:"subcategories" validate:"required"`
}

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
