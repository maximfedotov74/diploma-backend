package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateActionDto struct {
	EndDate     time.Time    `json:"end_date" validate:"required,gt"`
	Title       string       `json:"title" validate:"required,min=5"`
	ImgPath     *string      `json:"img_path" validate:"omitempty"`
	Description *string      `json:"description" validate:"omitempty,min=10"`
	Gender      ActionGender `json:"gender" validate:"required,actionGenderEnumValidation"`
}

type AddModelToActionDto struct {
	ProductModelId int    `json:"product_model_id" validate:"required,min=1"`
	ActionId       string `json:"action_id" validate:"required,uuid"`
}

type ActionGender string

const (
	MEN      ActionGender = "men"
	WOMEN    ActionGender = "women"
	CHILDREN ActionGender = "children"
	EVERYONE ActionGender = "everyone"
)

func ActionGenderEnumValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(MEN), string(WOMEN), string(CHILDREN), string(EVERYONE):
		return true
	}
	return false
}

type Action struct {
	Id          string       `json:"id" validate:"required"`
	CreatedAt   time.Time    `json:"created_at" validate:"required"`
	UpdatedAt   time.Time    `json:"updated_at" validate:"required"`
	EndDate     time.Time    `json:"end_date" validate:"required"`
	Title       string       `json:"title" validate:"required,min=5"`
	IsActivated bool         `json:"is_activated" validate:"required"`
	ImgPath     *string      `json:"img_path" validate:"omitempty"`
	Gender      ActionGender `json:"gender" validate:"required"`
	Description *string      `json:"description" validate:"omitempty,min=10"`
}

type ActionModel struct {
	ProductId     int           `json:"product_id" example:"1" validate:"required"`
	Title         string        `json:"product_title" example:"Ботинки" validate:"required"`
	Slug          string        `json:"product_slug" example:"botinki" validate:"required"`
	Article       string        `json:"article" validate:"required,min=1"`
	ModelId       int           `json:"model_id" example:"1" validate:"required"`
	ActionModelId int           `json:"action_model_id" example:"1" validate:"required"`
	Price         int           `json:"model_price" example:"10000" validate:"required"`
	Discount      *byte         `json:"model_discount" example:"15"`
	MainImagePath string        `json:"model_main_image_path" example:"/static/category/test.webp" validate:"required"`
	Brand         Brand         `json:"brand" validate:"required"`
	Category      CategoryModel `json:"category" validate:"required"`
}

type UpdateActionDto struct {
	EndDate     *time.Time `json:"end_date" validate:"omitempty,gt"`
	Title       *string    `json:"title" validate:"omitempty,min=5"`
	ImgPath     *string    `json:"img_path" validate:"omitempty"`
	Description *string    `json:"description" validate:"omitempty,min=10"`
	IsActivated *bool      `json:"is_activated" validate:"omitempty"`
}
