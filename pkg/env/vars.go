package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadVars() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return nil
}
