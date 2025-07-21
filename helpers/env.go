package helpers

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func InitEnv() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func LoadEnv(variableName string) (string, error) {
	val := os.Getenv(variableName)
	if val == "" {
		return "", fmt.Errorf("environment variable %s not set", variableName)
	}
	return val, nil
}
