package main

import (
	"log"
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

	// Init JWT secrets
	auth.Init(cfg.JWT.Secret, cfg.JWT.RefreshSecret)

	// Tạo GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			DB: db,
		},
	}))

	// Tạo router Gin
	r := gin.Default()

	// Public route: Playground
	r.GET("/graphql", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))

	// Protected route: chỉ Manager mới được vào
	r.POST("/query", auth.AuthMiddleware(), gin.WrapH(srv))

	// Start server
	log.Println("Server started at http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
