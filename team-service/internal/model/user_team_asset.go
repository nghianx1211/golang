package model

import "time"

type User struct {
    UserID       string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Username     string `gorm:"size:50;not null"`
    Email        string `gorm:"size:100;unique;not null"`
    Role         string `gorm:"size:20;not null"`
    PasswordHash string `gorm:"not null"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type Team struct {
    TeamID    string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    TeamName  string `gorm:"size:100;not null"`
    Members   []User `gorm:"many2many:team_members"`
    Managers  []User `gorm:"many2many:team_managers"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Folder struct {
    FolderID  string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name      string `gorm:"size:100;not null"`
    OwnerID   string `gorm:"type:uuid;not null"`
    Owner     User   // GORM tự map theo OwnerID → User.UserID
    Notes     []Note
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Note struct {
    NoteID    string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Title     string `gorm:"size:100"`
    Body      string
    FolderID  string `gorm:"type:uuid;not null"`
    Folder    Folder
    OwnerID   string `gorm:"type:uuid;not null"`
    Owner     User
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type FolderShare struct {
    FolderID  string `gorm:"type:uuid;primaryKey"`
    UserID    string `gorm:"type:uuid;primaryKey"`
    Access    string `gorm:"size:10;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type NoteShare struct {
    NoteID    string `gorm:"type:uuid;primaryKey"`
    UserID    string `gorm:"type:uuid;primaryKey"`
    Access    string `gorm:"size:10;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
