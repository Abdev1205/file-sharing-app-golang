package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file or no env file found")
	}
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
