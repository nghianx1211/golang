package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Folder model
type Folder struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Notes       []Note         `json:"notes,omitempty" gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE"`
	SharedWith  []FolderShare  `json:"shared_with,omitempty" gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE"`
}

// Note model
type Note struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"`
	FolderID  uuid.UUID `json:"folder_id" gorm:"type:uuid;not null;index"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Folder     Folder      `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
	SharedWith []NoteShare `json:"shared_with,omitempty" gorm:"foreignKey:NoteID;constraint:OnDelete:CASCADE"`
}

// FolderShare model for sharing folders
type FolderShare struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FolderID   uuid.UUID `json:"folder_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Permission string    `json:"permission" gorm:"not null;check:permission IN ('read', 'write')"`
	SharedBy   uuid.UUID `json:"shared_by" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
	
	// Relationships
	Folder Folder `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
}

// NoteShare model for sharing individual notes
type NoteShare struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NoteID     uuid.UUID `json:"note_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Permission string    `json:"permission" gorm:"not null;check:permission IN ('read', 'write')"`
	SharedBy   uuid.UUID `json:"shared_by" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
	
	// Relationships
	Note Note `json:"note,omitempty" gorm:"foreignKey:NoteID"`
}

// Request/Response DTOs
type CreateFolderRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateFolderRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=255"`
	Content string `json:"content"`
}

type ShareRequest struct {
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	Permission string    `json:"permission" binding:"required,oneof=read write"`
}

// Response models for better JSON output
type FolderResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Notes       []NoteResponse `json:"notes,omitempty"`
	SharedWith  []ShareInfo    `json:"shared_with,omitempty"`
}

type NoteResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	FolderID  uuid.UUID `json:"folder_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SharedWith []ShareInfo `json:"shared_with,omitempty"`
}

type ShareInfo struct {
	UserID     uuid.UUID `json:"user_id"`
	Permission string    `json:"permission"`
	SharedBy   uuid.UUID `json:"shared_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type AssetResponse struct {
	Folders []FolderResponse `json:"folders"`
	Notes   []NoteResponse   `json:"notes"`
}

// User info from user service
type UserInfo struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}