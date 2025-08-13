package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"asset-service/config"
	"asset-service/internal/database"
	"asset-service/internal/handler"
	"asset-service/internal/middleware"
	"asset-service/internal/service"
	route "asset-service/pkg/router"
)

func main() {
	// Cấu hình format log
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[ERROR] Failed to load config: %v", err)
	}

	// Kết nối DB & migrate
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect database: %v", err)
	}
	database.Migrate(db)
	log.Println("[INFO] Database connected and migrated successfully")

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
		AllowOrigins:     []string{"*"}, // TODO: Điều chỉnh cho production
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

	// Setup routes
	route.SetupAssetRoutes(r, assetHandler, authMiddleware)
	log.Println("[INFO] Routes configured successfully")

	// Get port from environment or use default
	port := getEnv("PORT", "8082")

	log.Printf("[INFO] Asset Management API starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
