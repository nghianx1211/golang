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
