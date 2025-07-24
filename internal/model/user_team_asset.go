package model

import (
    "time"
)

type User struct {
    UserID          string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Username        string `gorm:"size:50;not null"`
    Email           string `gorm:"size:100;unique;not null"`
    Role            string `gorm:"size:20;not null"`
    PasswordHash    string `gorm:"not null"`
}

type Team struct {
    TeamID    string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    TeamName  string `gorm:"size:100;not null"`
    Members   []User `gorm:"many2many:team_members;"`
    Managers  []User `gorm:"many2many:team_managers;"`
}

type Folder struct {
    FolderID string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name     string `gorm:"size:100;not null"`
    OwnerID  string `gorm:"type:uuid"`
    Owner    User
    Notes    []Note
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Note struct {
    NoteID   string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Title    string `gorm:"size:100"`
    Body     string
    FolderID string `gorm:"type:uuid"`
    Folder   Folder
    OwnerID  string `gorm:"type:uuid"`
    Owner    User
}

type FolderShare struct {
    FolderID string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID   string `gorm:"type:uuid"`
    Access   string `gorm:"size:10;not null"`
}

type NoteShare struct {
    NoteID string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID string `gorm:"type:uuid"`
    Access string `gorm:"size:10;not null"`
}