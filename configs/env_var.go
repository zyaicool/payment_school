package configs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	env := os.Getenv("GO_ENV")
	envFile := ".env"

	if env == "development" {
		envFile = ".env.dev"
	}

	envPath := filepath.Clean(envFile)
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading environment file: ", err)
	}
}
