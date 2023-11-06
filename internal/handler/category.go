package handler

import (
	"log"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

func (h *Handler) initCategoryRoutes(router fiber.Router) {
	category := router.Group("/category")
	{
		category.Post("/create-type", h.createCategoryType)
		category.Post("/create-category", h.createCategory)
		// category.Get("/type/:title")
		category.Get("/:title", h.getCategoryBytitle)
	}

}

// @Summary Create category type
// @Description Create category type with dto
// @Tags category
// @Accept json
// @Produce json
// @Param dto body model.CreateCategoryTypeDto true "Create category-type dto"
// @Router /api/category/create-type [post]
// @Success 201
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) createCategoryType(ctx *fiber.Ctx) error {
	dto := model.CreateCategoryTypeDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := lib.NewErr(messages.INVALID_BODY, 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := lib.ValidationMessages(error_messages)
		validError := lib.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.services.CategoryService.CreateCategoryType(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(201)
}

// @Summary Create category
// @Description Create category with dto
// @Tags category
// @Accept json
// @Produce json
// @Param dto body model.CreateCategoryDto true "Create category dto"
// @Router /api/category/create-category [post]
// @Success 201
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) createCategory(ctx *fiber.Ctx) error {
	dto := model.CreateCategoryDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := lib.NewErr(messages.INVALID_BODY, 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	log.Println(dto)

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := lib.ValidationMessages(error_messages)
		validError := lib.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.services.CategoryService.CreateCategory(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(201)
}

func (h *Handler) getCategoryBytitle(ctx *fiber.Ctx) error {
	undecodedTitle := ctx.Params("title")
	decodedTitle, err := url.PathUnescape(undecodedTitle)
	if err != nil {
		return ctx.SendStatus(400)
	}

	category, appErr := h.services.CategoryService.FindCategoryByTitle(decodedTitle)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(200).JSON(category)

}
