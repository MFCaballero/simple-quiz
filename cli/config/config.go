package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BackendURL string `required:"true"`
}

func LoadConfig() Config {
	var config Config

	envPath := os.Getenv("DOTENV")
	if envPath == "" {
		envPath = "cli/.env"
	}

	if err := godotenv.Load(envPath); err != nil {
		panic(err.Error())
	}

	if err := envconfig.Process("", &config); err != nil {
		panic(err.Error())
	}
	return config
}
