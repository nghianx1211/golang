package database

import (
	"gorm.io/gorm"
	"golang/internal/model"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Team{},
		&model.Folder{},
		&model.Note{},
		&model.FolderShare{},
		&model.NoteShare{},
	)
}