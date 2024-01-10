package config

import (
	"backend/db/migrate"
	"backend/db/seeder"
	"backend/internal/delivery/http"
	"backend/internal/delivery/http/middleware"
	"backend/internal/delivery/http/route"
	"backend/internal/repository"
	"backend/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App      *echo.Echo
	DB       *gorm.DB
	Redis    *redis.Client
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository()
	postRepository := repository.NewPostRepository()

	// setup usecases
	postUseCase := usecase.NewPostUseCase(config.DB, config.Redis, config.Validate, postRepository, userRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Redis, config.Validate, userRepository, config.Config)

	// setup controller
	postController := http.NewPostController(postUseCase)
	userController := http.NewUserController(userUseCase)

	// setup middleware
	authMiddleware := middleware.AuthMiddleware(config.Config, config.Redis)

	// setup route
	routeConfig := route.RouteConfig{
		App:            config.App,
		PostController: postController,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()

	// migrate the database
	var err error
	err = migrate.Drop(config.DB)
	err = migrate.Migrate(config.DB)
	users, err := seeder.SeedsUser(config.DB, userRepository)
	_, err = seeder.SeedsPost(config.DB, postRepository, &users)

	if err != nil {
		panic(err)
	}
}
