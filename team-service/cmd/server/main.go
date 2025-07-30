package main

import (
	"log"
	"os"

	"team-service/internal/handler"   
	"team-service/internal/model"     
	"team-service/config"     
	"team-service/internal/database"     
	route "team-service/pkg/router"     
	"team-service/internal/service"   

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	userServiceClient := service.NewUserServiceClient(getEnv("USER_SERVICE_URL", "http://localhost:8080"))
	teamService := service.NewTeamService(db, userServiceClient)

	// Initialize handler
	teamHandler := handler.NewTeamHandler(teamService)

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware if needed
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup router
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	route.SetupTeamRouter(router, teamHandler, jwtSecret)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := getEnv("PORT", "8081")
	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Team{},
		&model.Manager{},
		&model.Member{},
	)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
