package main

import (
	"os"
	"team-service/config"
	"team-service/internal/database"
	"team-service/internal/handler"
	"team-service/internal/messaging"
	"team-service/internal/model"
	"team-service/internal/service"
	route "team-service/pkg/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var log = logrus.New()

func main() {
	// Setup log format JSON
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to load config")
	}

	// Connect DB & migrate
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect database")
	}
	database.Migrate(db)
	log.Info("Database connected and migrated successfully")

	redisClient := messaging.NewRedisClient(getEnv("REDIS_ADDR", "localhost:6379"))
	teamProducer := messaging.NewKafkaProducer(
		getEnv("KAFKA_BROKER", "localhost:9092"),
		"team.activity",
	)
	// Initialize services
	userServiceClient := service.NewUserServiceClient(getEnv("USER_SERVICE_URL", "http://localhost:8080"))
	teamService := service.NewTeamService(db, userServiceClient, teamProducer, redisClient)

	// Initialize handler
	teamHandler := handler.NewTeamHandler(teamService)

	// Initialize router
	router := gin.Default()

	// CORS middleware
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

	// Setup routes
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	route.SetupTeamRouter(router, teamHandler, jwtSecret)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		log.Info("Health check endpoint called")
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := getEnv("PORT", "8081")
	log.WithField("port", port).Info("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.WithError(err).Fatal("Server failed to start")
	}
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Team{},
		&model.Manager{},
		&model.Member{},
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
