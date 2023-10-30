package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/service"
)

type Handler struct {
	services *service.Services
	cfg      *cfg.Config
}

func New(services *service.Services, cfg *cfg.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) Init(cfg *cfg.Config, router fiber.Router) {
	h.initUsersRoutes(router)
	h.initRoleRoutes(router)
	h.initAuthRoutes(router)
}
