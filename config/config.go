package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBDriver     string
	DBConnection string
}

func NewConfig() *Config {

	envFile := ".env"

	injectedEnvFile := os.Getenv("ENV_FILE")
	if injectedEnvFile != "" {
		envFile = injectedEnvFile
	}

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println(err)
	}

	config := &Config{
		Port:         os.Getenv("PORT"),
		DBDriver:     os.Getenv("DB_DRIVER"),
		DBConnection: os.Getenv("DB_CONNECTION"),
	}

	if config.Port == "" {
		config.Port = "8000"
	}

	return config
}
