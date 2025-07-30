package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	
	"asset-service/internal/handler"
	"asset-service/internal/middleware"
	route "asset-service/pkg/router"
	"asset-service/internal/service"
	"asset-service/config"     
	"asset-service/internal/database" 
)

func main() {
	// Load cấu hình
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Kết nối DB & migrate
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	database.Migrate(db)

	// Initialize service
	userServiceClient := service.NewUserServiceClient()
	assetService := service.NewAssetService(db, userServiceClient)

	// Initialize handler
	assetHandler := handler.NewAssetHandler(assetService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userServiceClient)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "asset-management-api",
		})
	})

	// Setup route
	route.SetupAssetRoutes(r, assetHandler, authMiddleware)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Asset Management API starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}