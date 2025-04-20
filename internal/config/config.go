package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, errors.New("Error loading .env file")
	}

	// Load environment variables
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	if mysqlDatabase == "" {
		log.Fatal("MYSQL_DATABASE not set")
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		log.Fatal("MYSQL_USER not set")
	}

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		log.Fatal("MYSQL_PASSWORD not set")
	}

	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		log.Fatal("MYSQL_HOST not set")
	}

	return &Config{
		DBName:     mysqlDatabase,
		DBUser:     mysqlUser,
		DBPassword: mysqlPassword,
		DBHost:     mysqlHost,
	}, nil
}
