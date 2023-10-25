package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/maximfedotov74/fiber-psql/docs"
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/handler"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/internal/scheduler"
	"github.com/maximfedotov74/fiber-psql/internal/service"
	"github.com/maximfedotov74/fiber-psql/pkg/db"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @Title Fiber Golang Api
// @Version 1.0
// Description This is a simple REST API using go fiber and postgresql
// @Contact.name Maxim Fedotov
// @Contact.url https://github.com/maximfedotov74
func Start() {

	config := cfg.GetCfg()

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New(cfg.CorsCfg()))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	dbClient := db.NewClient(config.DatabaseUrl)
	//cacheManager := cache.NewCacheClient(config.RedisAddr, config.RedisPassword)

	repositories := repository.New(dbClient)
	services := service.New(service.Deps{Repos: repositories, Config: config})
	handler := handler.New(services, config)

	router := app.Group("/api")

	handler.Init(config, router)

	scheduler := scheduler.New(handler)
	scheduler.Start()

	PORT := config.Port

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	go func() {
		<-quit
		log.Info("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	log.Infof("Server started on PORT: %s", PORT)
	if err := app.Listen(fmt.Sprintf(":%s", PORT)); err != nil {
		log.Fatal(err)
	}

	log.Info("Cleaning")
	scheduler.Shutdown()
	dbClient.Close()
}
