package app

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/maximfedotov74/diploma-backend/docs"
	"github.com/maximfedotov74/diploma-backend/internal/config"
	"github.com/maximfedotov74/diploma-backend/internal/domain/handler"
	"github.com/maximfedotov74/diploma-backend/internal/domain/repository"
	"github.com/maximfedotov74/diploma-backend/internal/domain/service"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Start() {
	configuration := config.MustLoadConfig()

	fiberApp := fiber.New(fiber.Config{BodyLimit: 10 * 10 * 1024 * 1024})

	fiberApp.Use((logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	})))

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: configuration.ClientUrl,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowCredentials: true,
	}))

	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	router := fiberApp.Group("/api")

	initDeps(router)

	log.Infof("Swagger Api docs working on : %s", "/swagger")
	log.Infof("Server started on PORT: %s", configuration.Port)
	go func() {
		if err := fiberApp.Listen(fmt.Sprintf(":%s", configuration.Port)); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit
	log.Info("Gracefully shutting down...")
	log.Info("Cleaning")
	fiberApp.Shutdown()
	log.Info("Application shutdown successfully!")
}

func initDeps(router fiber.Router) {
	deliveryRepository := repository.NewDeliveryRepository()
	deliveryService := service.NewDeliveryService(deliveryRepository)
	deliveryHandler := handler.NewDeliveryHandler(deliveryService)
	router.Post("/test", deliveryHandler.Create)
}
