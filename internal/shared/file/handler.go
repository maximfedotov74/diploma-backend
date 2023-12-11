package file

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type Service interface {
	Upload(dir string, header *multipart.FileHeader) (*string, error)
}

type FileHandler struct {
	service Service
	router  fiber.Router
}

func NewFileHandler(service Service, router fiber.Router) *FileHandler {
	return &FileHandler{service: service, router: router}
}

func (fh *FileHandler) InitRoutes() {
	fileRouter := fh.router.Group("file")
	{
		fileRouter.Post("/", fh.upload)
	}
}

func (fh *FileHandler) upload(ctx *fiber.Ctx) error {
	f, err := ctx.FormFile("file")

	dir := ctx.Query("dir", "default")

	if err != nil {
		return err
	}

	path, err := fh.service.Upload(dir, f)
	if err != nil {
		return err
	}
	return ctx.Status(exception.STATUS_CREATED).JSON(FileUploadResponse{Path: *path})
}
