package api

import (
	"ciphersafe/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures the application's routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Instantiate services
	userService := services.NewUserService(db)
	authService := services.NewAuthService(userService)

	// Instantiate handlers
	authHandler := NewAuthHandler(authService)
	projectHandler := NewProjectHandler(db)
	secretHandler := NewSecretHandler(db)

	// Public routes (auth)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Protected routes (main API)
	api := r.Group("/api")
	api.Use(AuthMiddleware())
	{
		// Project routes
		api.POST("/projects", projectHandler.CreateProject)
		api.GET("/projects", projectHandler.GetProjects)

		// Secret routes
		api.POST("/secrets", secretHandler.CreateSecret)
		api.GET("/projects/:projectID/secrets", secretHandler.GetSecretsForProject)
		api.DELETE("/secrets/:secretID", secretHandler.DeleteSecret)
	}
}
