package main

import (
	config "user-service/config"
	"user-service/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, _ := database.Connect(cfg.Database)
	database.Migrate(db)
}
