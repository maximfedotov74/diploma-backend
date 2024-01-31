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
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/repository"
	"github.com/maximfedotov74/diploma-backend/internal/domain/service"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
	"github.com/maximfedotov74/diploma-backend/internal/shared/mail"
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

	postgresClient := db.NewPostgresConnection(configuration.DatabaseUrl)

	router := fiberApp.Group("/api")

	initDeps(router, postgresClient, configuration)

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
	postgresClient.Close()
	log.Info("Application shutdown successfully!")
}

func initDeps(router fiber.Router, postgresClient db.PostgresClient, config *config.Config) {

	jwtService := jwt.NewJwtService(jwt.JwtConfig{
		RefreshTokenExp:    config.RefreshTokenExp,
		AccessTokenExp:     config.AccessTokenExp,
		RefreshTokenSecret: config.RefreshTokenSecret,
		AccessTokenSecret:  config.AccessTokenSecret,
	})

	mailService := mail.NewMailService(mail.MailConfig{SmtpKey: config.SmtpKey, SenderEmail: config.SmtpMail, SmtpHost: config.SmtpHost, SmtpPort: config.SmtpPort, AppLink: config.AppLink})

	sessionRepo := repository.NewSessionRepository(postgresClient)
	sessionService := service.NewSessionService(sessionRepo, jwtService)

	roleRepo := repository.NewRoleRepository(postgresClient)
	userRepo := repository.NewUserRepository(postgresClient, roleRepo)
	brandRepo := repository.NewBrandRepository(postgresClient)
	categoryRepo := repository.NewCategoryRepository(postgresClient)
	deliveryRepo := repository.NewDeliveryRepository(postgresClient)
	optionRepo := repository.NewOptionRepository(postgresClient)
	productRepo := repository.NewProductRepository(postgresClient)
	feedbackRepo := repository.NewFeedbackRepository(postgresClient)

	roleService := service.NewRoleService(roleRepo)
	userService := service.NewUserService(userRepo)
	brandService := service.NewBrandService(brandRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	optionService := service.NewOptionService(optionRepo)
	productService := service.NewProductService(productRepo, brandService, categoryService)
	feedbackService := service.NewFeedbackService(feedbackRepo)

	authMiddleware := middleware.CreateAuthMiddleware(sessionService, userService)
	//roleMiddleware := middleware.CreateRoleMiddleware()

	authService := service.NewAuthService(userService, sessionService, mailService)
	roleHandler := handler.NewRoleHandler(roleService, router)
	userHandler := handler.NewUserHandler(userService, router)
	authHandler := handler.NewAuthHandler(authService, router, authMiddleware)
	brandHandler := handler.NewBrandHandler(brandService, router, authMiddleware)
	categoryHandler := handler.NewCategoryHandler(categoryService, router, authMiddleware)
	deliveryHandler := handler.NewDeliveryHandler(deliveryRepo, router, authMiddleware)
	optionHandler := handler.NewOptionHandler(optionService, router, authMiddleware)
	productHandler := handler.NewProductHandler(productService, router, authMiddleware)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService, router, authMiddleware)

	roleHandler.InitRoutes()
	userHandler.InitRoutes()
	authHandler.InitRoutes()
	brandHandler.InitRoutes()
	categoryHandler.InitRoutes()
	deliveryHandler.InitRoutes()
	optionHandler.InitRoutes()
	productHandler.InitRoutes()
	feedbackHandler.InitRoutes()
}
