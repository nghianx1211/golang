package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Team represents a team entity
type Team struct {
	TeamID    string    `json:"teamId" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TeamName  string    `json:"teamName" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Managers  []Manager `json:"managers" gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
	Members   []Member  `json:"members" gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
}

// Manager represents a team manager
type Manager struct {
	ID          uint   `json:"-" gorm:"primaryKey"`
	TeamID      string `json:"-" gorm:"type:uuid;not null"`
	ManagerID   string `json:"managerId" gorm:"type:uuid;not null"`
	ManagerName string `json:"managerName" gorm:"not null"`
	IsMain      bool   `json:"isMain" gorm:"default:false"` // Only main manager can add/remove other managers
	CreatedAt   time.Time `json:"createdAt"`
}

// Member represents a team member
type Member struct {
	ID         uint   `json:"-" gorm:"primaryKey"`
	TeamID     string `json:"-" gorm:"type:uuid;not null"`
	MemberID   string `json:"memberId" gorm:"type:uuid;not null"`
	MemberName string `json:"memberName" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt"`
}

// BeforeCreate hook for Team
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.TeamID == "" {
		t.TeamID = uuid.New().String()
	}
	return nil
}

// CreateTeamRequest represents the request body for creating a team
type CreateTeamRequest struct {
	TeamName string `json:"teamName" binding:"required"`
	Managers []struct {
		ManagerID   string `json:"managerId" binding:"required"`
		ManagerName string `json:"managerName" binding:"required"`
	} `json:"managers" binding:"required,min=1"`
	Members []struct {
		MemberID   string `json:"memberId" binding:"required"`
		MemberName string `json:"memberName" binding:"required"`
	} `json:"members"`
}

// AddMemberRequest represents the request body for adding a member
type AddMemberRequest struct {
	MemberID   string `json:"memberId" binding:"required"`
	MemberName string `json:"memberName" binding:"required"`
}

// AddManagerRequest represents the request body for adding a manager
type AddManagerRequest struct {
	ManagerID   string `json:"managerId" binding:"required"`
	ManagerName string `json:"managerName" binding:"required"`
}

// UserServiceResponse represents user data from user service
type UserServiceResponse struct {
	Data struct {
		FetchUsers []struct {
			UserID   string `json:"userID"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		} `json:"fetchUsers"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}