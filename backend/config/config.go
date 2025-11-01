package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MasterEncryptionKey []byte
	JWTSecretKey        []byte
	DatabaseURL         string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	key := os.Getenv("MASTER_ENCRYPTION_KEY")
	if key == "" {
		log.Fatal("MASTER_ENCRYPTION_KEY is not set")
	}

	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	AppConfig = &Config{
		MasterEncryptionKey: []byte(key), // Store as bytes
		JWTSecretKey:        []byte(jwtKey),
		DatabaseURL:         dbURL,
	}
}
