package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/maximfedotov74/fiber-psql/docs"
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/domain/auth"
	"github.com/maximfedotov74/fiber-psql/internal/domain/category"
	"github.com/maximfedotov74/fiber-psql/internal/domain/role"
	"github.com/maximfedotov74/fiber-psql/internal/domain/session"
	"github.com/maximfedotov74/fiber-psql/internal/domain/user"
	"github.com/maximfedotov74/fiber-psql/internal/guards"
	"github.com/maximfedotov74/fiber-psql/internal/shared/cache"
	"github.com/maximfedotov74/fiber-psql/internal/shared/db"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/mail"
	"github.com/maximfedotov74/fiber-psql/internal/shared/password"
	"github.com/maximfedotov74/fiber-psql/internal/shared/scheduler"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Application struct{}

type Test struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

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

	fiberInstance := fiber.New()

	fiberInstance.Use((logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	})))

	fiberInstance.Get("/swagger/*", fiberSwagger.WrapHandler)

	dbService := db.NewDbService(cfg.DatabaseUrl)

	router := fiberInstance.Group("/api")

	cron := gocron.NewScheduler(time.UTC)

	schdulerService := scheduler.New(cron)
	schdulerService.Start()

	cacheContext, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cacheService := cache.NewCacheService(cfg.RedisAddr, cfg.RedisPassword, cacheContext)

	max := Test{Id: 1, Name: "Maxim", Age: 19}
	err := cacheService.Set(strconv.Itoa(max.Id), max, time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	res := Test{}
	err = cacheService.Get("1", &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(res)

	app.initializeDependencies(dbService, cfg, router, cacheService)

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
	cacheService.Shutdown()
	schdulerService.Shutdown()
	dbService.Close()
	log.Info("Application shutdown successfully!")

}

func (app *Application) initializeDependencies(dbService *pgxpool.Pool, cfg *cfg.Config,
	router fiber.Router, cacheService *cache.CacheService) {
	jwtSerivce := jwt.NewJwtService(jwt.JwtConfig{RefreshTokenExp: cfg.RefreshTokenExp, AccessTokenExp: cfg.AccessTokenExp, RefreshTokenSecret: cfg.RefreshTokenSecret, AccessTokenSecret: cfg.AccessTokenSecret})

	mailService := mail.NewMailService(mail.MailConfig{SmtpKey: cfg.SmtpKey, SenderEmail: cfg.SmtpMail, SmtpHost: cfg.SmtpHost, SmtpPort: cfg.SmtpPort, AppLink: cfg.AppLink})

	passwordService := password.NewPasswordService()

	roleRepository := role.NewRoleRepository(dbService)
	categoryRepository := category.NewCategoryRepository(dbService)
	userRepository := user.NewUserRepository(dbService, roleRepository)
	sessionRepository := session.NewSessionRepository(dbService)

	roleService := role.NewRoleService(roleRepository)
	categoryService := category.NewCategoryService(categoryRepository)
	sessionService := session.NewSessionService(sessionRepository, jwtSerivce)

	userService := user.NewUserService(userRepository, sessionService, mailService, passwordService)
	authService := auth.NewAuthService(userService, sessionService, passwordService, mailService)

	authGuard := guards.NewAuthGuard(sessionService, userService)
	roleGuard := guards.NewRoleGuard()

	userHandler := user.NewUserHandler(userService, router, authGuard.CheckAuth)
	roleHandler := role.NewRoleHandler(roleService, authGuard.CheckAuth, roleGuard, router)
	categoryHandler := category.NewCategoryHandler(categoryService, router, authGuard.CheckAuth)
	authHandler := auth.NewAuthHandler(authService, router, authGuard.CheckAuth)

	userHandler.InitRoutes()
	roleHandler.InitRoutes()
	categoryHandler.InitRoutes()
	authHandler.InitRoutes()
}
