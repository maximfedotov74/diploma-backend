package handler

import (
	"context"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type fileClient interface {
	Upload(ctx context.Context, h *multipart.FileHeader) (*model.UploadResponse, error)
}

type FileHandler struct {
	client         fileClient
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewFileHandler(client fileClient, router fiber.Router, authMiddleware middleware.AuthMiddleware) *FileHandler {
	return &FileHandler{client: client, router: router, authMiddleware: authMiddleware}
}

func (h *FileHandler) InitRoutes() {
	fileRouter := h.router.Group("file")
	{
		fileRouter.Post("/", h.upload)
	}
}

// @Summary Upload file
// @Description Upload file
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Router /api/file/ [post]
// @Success 201 {object} model.UploadResponse
// @Failure 400 {object} fall.AppErr
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *FileHandler) upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		ex := fall.ServerError(err.Error())
		return ctx.Status(ex.Status()).JSON(ex)
	}
	res, err := h.client.Upload(ctx.Context(), file)
	if err != nil {
		ex := fall.ServerError(err.Error())
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(fall.STATUS_CREATED).JSON(res)
}
