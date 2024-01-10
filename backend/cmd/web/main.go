package main

import (
	"backend/internal/config"
	"log"

	"github.com/go-playground/validator/v10"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewEcho()
	db := config.NewDatabase(viperConfig)
	redis := config.NewRedisClient(viperConfig)
	validate := validator.New()

	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		DB:       db,
		Redis:    redis,
		Validate: validate,
		Config:   viperConfig,
	})

	port := viperConfig.GetString("web.port")
	log.Fatal(app.Start(":" + port))
}
