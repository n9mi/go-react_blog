package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper) *redis.Client {
	// Getting config from viper
	address := viper.GetString("redis.address")
	port := viper.GetInt("redis.port")
	db := viper.GetInt("redis.db")
	password := viper.GetString("redis.password")

	// Create new client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		DB:       db,
		Password: password,
	})
	// Delete all redis value
	client.FlushDB(context.Background())

	return client
}
