package main

import (
	"ciphersafe/api"
	"ciphersafe/config"
	"ciphersafe/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load config from .env
	config.LoadConfig()

	// 2. Connect to Database
	db, err := gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 3. Auto-migrate the schema
	log.Println("Migrating database...")
	db.AutoMigrate(&models.User{}, &models.Project{}, &models.Secret{})

	// 4. Set up Gin router
	r := gin.Default()

	// 5. Setup routes
	api.SetupRoutes(r, db)

	// 6. Start server
	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
