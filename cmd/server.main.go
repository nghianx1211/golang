package main

import (
    "golang/internal/database"
    "golang/pkg/router"
	"golang/configs"
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
