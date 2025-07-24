package database

import (
	"github.com/nghianx1211/golang/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
    // Enable UUID extension first
    if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
        return err
    }

    // Correct order: parent tables before child tables
    return db.AutoMigrate(
        &model.User{},
        &model.Team{},
        &model.Folder{},
        &model.Note{},
        &model.FolderShare{},
        &model.NoteShare{},
    )
}