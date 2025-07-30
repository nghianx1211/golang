package database

import (
	"team-service/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
    // Enable UUID extension first
    if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
        return err
    }

    // Correct order: parent tables before child tables
    return db.AutoMigrate(
		&model.Team{},
		&model.Manager{},
		&model.Member{},
    )
}