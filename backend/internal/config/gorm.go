package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase returns new gorm database instance
func NewDatabase(viper *viper.Viper) *gorm.DB {
	// Getting config from viper
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	name := viper.GetString("database.name")
	port := viper.GetInt("database.port")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifetimeConnection := viper.GetInt("database.pool.lifetime")

	// Create dsn from config
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		host,
		username,
		password,
		name,
		port)

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// Getting database instance
	connection, err := db.DB()

	if err != nil {
		panic(err)
	}

	// Set connection configuration
	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxIdleConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifetimeConnection))

	return db
}
