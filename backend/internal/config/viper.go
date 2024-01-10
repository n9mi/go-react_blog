package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper
func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")    // Define config.json path to working directory
	config.AddConfigPath("./../") // Define config.json path to working directory

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return config
}
