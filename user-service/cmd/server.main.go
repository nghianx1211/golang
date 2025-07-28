package main

import (
	config "github.com/nghianx1211/golang/configs"
	"github.com/nghianx1211/golang/internal/database"
	"github.com/nghianx1211/golang/pkg/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db := database.InitDB(cfg.DB_DSN)
	database.Migrate(db)

	r := router.Setup(db)
	r.Run(":8080")
}
