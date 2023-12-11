package app

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/maximfedotov74/fiber-psql/docs"
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/domain/auth"
	"github.com/maximfedotov74/fiber-psql/internal/domain/brand"
	"github.com/maximfedotov74/fiber-psql/internal/domain/category"
	"github.com/maximfedotov74/fiber-psql/internal/domain/feedback"
	"github.com/maximfedotov74/fiber-psql/internal/domain/option"
	"github.com/maximfedotov74/fiber-psql/internal/domain/product"
	"github.com/maximfedotov74/fiber-psql/internal/domain/role"
	"github.com/maximfedotov74/fiber-psql/internal/domain/session"
	"github.com/maximfedotov74/fiber-psql/internal/domain/user"
	"github.com/maximfedotov74/fiber-psql/internal/domain/wish"
	"github.com/maximfedotov74/fiber-psql/internal/guards"
	"github.com/maximfedotov74/fiber-psql/internal/shared/db"
	"github.com/maximfedotov74/fiber-psql/internal/shared/file"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/mail"
	"github.com/maximfedotov74/fiber-psql/internal/shared/password"
	"github.com/maximfedotov74/fiber-psql/internal/shared/scheduler"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Application struct{}

func NewApplication() *Application {
	return &Application{}
}

// @Title Fiber Golang Api
// @Version 1.0
// Description This is a simple REST API using go fiber and postgresql
// @Contact.name Maxim Fedotov
// @Contact.url https://github.com/maximfedotov74
func (app *Application) Start() {
	cfg := cfg.MustGetCfg()

	fiberInstance := fiber.New(fiber.Config{BodyLimit: 10 * 10 * 1024 * 1024})

	fiberInstance.Use((logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	})))

	fiberInstance.Use(cors.New(cors.Config{
		AllowOrigins: cfg.ClientUrl,
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

	fiberInstance.Get("/swagger/*", fiberSwagger.WrapHandler)

	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	staticPath := path.Join(dir, cfg.StaticPath)

	_, err = os.Stat(staticPath)

	if os.IsNotExist(err) {
		os.Mkdir(staticPath, 0700)
	}

	fileService := file.NewFileService(staticPath)

	fiberInstance.Static("/static", staticPath, fiber.Static{
		ByteRange: true,
	})

	dbService := db.NewDbService(cfg.DatabaseUrl)

	router := fiberInstance.Group("/api")

	cron := gocron.NewScheduler(time.UTC)

	schdulerService := scheduler.New(cron)
	schdulerService.Start()

	app.initializeDependencies(dbService, cfg, router, fileService)

	PORT := cfg.Port

	log.Infof("Swagger Api docs working on : %s", "/swagger")
	log.Infof("Server started on PORT: %s", PORT)
	go func() {
		if err := fiberInstance.Listen(fmt.Sprintf(":%s", PORT)); err != nil {
			log.Fatal(err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit
	log.Info("Gracefully shutting down...")

	log.Info("Cleaning")
	_ = fiberInstance.Shutdown()
	schdulerService.Shutdown()
	dbService.Close()
	log.Info("Application shutdown successfully!")

}

func (app *Application) initializeDependencies(dbService *pgxpool.Pool, cfg *cfg.Config,
	router fiber.Router, fileService *file.FileService) {

	jwtSerivce := jwt.NewJwtService(jwt.JwtConfig{RefreshTokenExp: cfg.RefreshTokenExp, AccessTokenExp: cfg.AccessTokenExp, RefreshTokenSecret: cfg.RefreshTokenSecret, AccessTokenSecret: cfg.AccessTokenSecret})

	mailService := mail.NewMailService(mail.MailConfig{SmtpKey: cfg.SmtpKey, SenderEmail: cfg.SmtpMail, SmtpHost: cfg.SmtpHost, SmtpPort: cfg.SmtpPort, AppLink: cfg.AppLink})

	passwordService := password.NewPasswordService()

	roleRepository := role.NewRoleRepository(dbService)
	categoryRepository := category.NewCategoryRepository(dbService)
	userRepository := user.NewUserRepository(dbService, roleRepository)
	sessionRepository := session.NewSessionRepository(dbService)
	brandRepository := brand.NewBrandRepository(dbService)
	productRepository := product.NewProductRepository(dbService)
	optionRepository := option.NewOptionRepository(dbService)
	feedbackRepository := feedback.NewFeedbackRepository(dbService)
	wishRepository := wish.NewWishRepository(dbService)

	roleService := role.NewRoleService(roleRepository)
	categoryService := category.NewCategoryService(categoryRepository)
	sessionService := session.NewSessionService(sessionRepository, jwtSerivce)
	brandService := brand.NewBrandService(brandRepository)
	productService := product.NewProductService(productRepository, categoryService, brandService)
	optionService := option.NewOptionService(optionRepository)
	feedbackService := feedback.NewFeedbackService(feedbackRepository)
	wishService := wish.NewWishService(wishRepository)

	userService := user.NewUserService(userRepository, sessionService, mailService, passwordService)
	authService := auth.NewAuthService(userService, sessionService, passwordService, mailService)

	authGuard := guards.NewAuthGuard(sessionService, userService)
	roleGuard := guards.NewRoleGuard()

	userHandler := user.NewUserHandler(userService, router, authGuard.CheckAuth)
	roleHandler := role.NewRoleHandler(roleService, authGuard.CheckAuth, roleGuard.CheckRoles, router)
	categoryHandler := category.NewCategoryHandler(categoryService, router, authGuard.CheckAuth)
	authHandler := auth.NewAuthHandler(authService, router, authGuard.CheckAuth)
	brandHandler := brand.NewBrandHandler(brandService, router)
	productHandler := product.NewProductHandler(productService, router, authGuard.CheckAuth, roleGuard.CheckRoles)
	optionHandler := option.NewOptionHandler(optionService, router, authGuard.CheckAuth)
	feedbackHandler := feedback.NewFeedbackHandler(feedbackService, router, authGuard.CheckAuth, roleGuard.CheckRoles)
	wishHandler := wish.NewWishHandler(wishService, router, authGuard.CheckAuth)
	fileHandler := file.NewFileHandler(fileService, router)

	userHandler.InitRoutes()
	roleHandler.InitRoutes()
	categoryHandler.InitRoutes()
	authHandler.InitRoutes()
	brandHandler.InitRoutes()
	productHandler.InitRoutes()
	optionHandler.InitRoutes()
	feedbackHandler.InitRoutes()
	wishHandler.InitRoutes()
	fileHandler.InitRoutes()
}
