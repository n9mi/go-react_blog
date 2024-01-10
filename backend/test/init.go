package test

import (
	"backend/internal/config"
	"backend/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	app         *echo.Echo
	db          *gorm.DB
	redisClient *redis.Client
	validate    *validator.Validate
	viperConfig *viper.Viper
)

var (
	validToken string
	authData   *model.UserAuthData
)

type TestSchema map[string]interface{}
type TestResponse[T any] struct {
	Code     int      `json:"code"`
	Status   string   `json:"status"`
	Data     T        `json:"data"`
	Messages []string `json:"messages"`
}

func init() {
	viperConfig = config.NewViper()
	app = config.NewEcho()
	db = config.NewDatabase(viperConfig)
	redisClient = config.NewRedisClient(viperConfig)
	validate = validator.New()

	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		DB:       db,
		Redis:    redisClient,
		Validate: validate,
		Config:   viperConfig,
	})
}
