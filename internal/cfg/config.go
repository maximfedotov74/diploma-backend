package cfg

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseUrl        string
	RefreshTokenExp    int
	AccessTokenExp     int
	RefreshTokenSecret string
	AccessTokenSecret  string
	SmtpKey            string
	SmtpHost           string
	SmtpPort           int
	SmtpMail           string
	AppLink            string
	RedisAddr          string
	RedisPassword      string
	DaDataKey          string
	StaticPath         string
	ClientUrl          string
}

func createPanicMessage(key string) string {
	return fmt.Sprintf("Key \"%s\" does not exist inside the env file", key)
}

var config *Config

var once sync.Once

func MustGetCfg() *Config {
	once.Do(func() {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal(".Env file not found!")
		}
		config = &Config{
			Port:               getEnv("PORT"),
			DatabaseUrl:        getEnv("DB_URL"),
			RefreshTokenExp:    getEnvAsInt("REFRESH_TOKEN_EXP"),
			AccessTokenExp:     getEnvAsInt("ACCESS_TOKEN_EXP"),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET"),
			AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET"),
			SmtpKey:            getEnv("SMTP_KEY"),
			SmtpPort:           getEnvAsInt("SMTP_PORT"),
			SmtpHost:           getEnv("SMTP_HOST"),
			SmtpMail:           getEnv("SMTP_MAIl"),
			AppLink:            getEnv("APP_LINK"),
			DaDataKey:          getEnv("DADATA_API_KEY"),
			StaticPath:         getEnv("STATIC_FILES_PATH"),
			ClientUrl:          getEnv("CLIENT_URL"),
		}
	})
	return config
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		err := createPanicMessage(key)
		panic(err)
	}

	return value
}

func getEnvAsInt(name string) int {
	valueStr := getEnv(name)
	value, err := strconv.Atoi(valueStr)

	if err != nil {
		panic("Error when coverting env variable to integer")
	}

	return value
}

func getEnvAsBool(name string) bool {
	valueStr := getEnv(name)
	value, err := strconv.ParseBool(valueStr)

	if err != nil {
		panic("Error when coverting env variable to boolean")
	}

	return value
}
