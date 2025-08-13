package main

import (
	"log/slog"
	"os"
	"user-service/config"
	"user-service/graph/generated"
	"user-service/graph/resolver"
	"user-service/internal/auth"
	"user-service/internal/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func main() {
	// Logger JSON cho Promtail
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	logger.Info("Hello Promtail", "service", "user-service")

	// Kết nối DB & migrate
	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect database", "error", err)
		os.Exit(1)
	}
	database.Migrate(db)
	logger.Info("Database connected and migrated", "service", "user-service")

	// Init JWT secrets
	auth.Init(cfg.JWT.Secret, cfg.JWT.RefreshSecret)
	logger.Info("JWT secrets initialized", "service", "user-service")

	// GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{DB: db},
	}))

	// Gin router
	r := gin.Default()
	r.GET("/graphql", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))
	r.POST("/query", auth.AuthMiddleware(), gin.WrapH(srv))

	// Start server
	logger.Info("Server started", "url", "http://localhost:8080", "service", "user-service")
	if err := r.Run(":8080"); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
