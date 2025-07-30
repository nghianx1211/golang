package service

import (
	"errors"
	"fmt"
	"strings"

	"team-service/internal/model"
	"gorm.io/gorm"
)

type TeamService struct {
	db                *gorm.DB
	userServiceClient *UserServiceClient
}

func NewTeamService(db *gorm.DB, userServiceClient *UserServiceClient) *TeamService {
	return &TeamService{
		db:                db,
		userServiceClient: userServiceClient,
	}
}

func (s *TeamService) CreateTeam(req *model.CreateTeamRequest, token string) (*model.Team, error) {
	// Validate all managers exist and have manager/admin role
	for _, manager := range req.Managers {
		_, err := s.userServiceClient.ValidateRole(manager.ManagerID, []string{"manager"}, token)
		if err != nil {
			return nil, fmt.Errorf("manager validation failed for %s: %v", manager.ManagerID, err)
		}
	}

	// Validate all members exist
	for _, member := range req.Members {
		_, err := s.userServiceClient.ValidateUser(member.MemberID, token)
		if err != nil {
			return nil, fmt.Errorf("member validation failed for %s: %v", member.MemberID, err)
		}
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create team
	team := &model.Team{
		TeamName: req.TeamName,
	}

	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create team: %v", err)
	}

	// Add managers
	for i, manager := range req.Managers {
		teamManager := model.Manager{
			TeamID:      team.TeamID,
			ManagerID:   manager.ManagerID,
			ManagerName: manager.ManagerName,
			IsMain:      i == 0, // First manager is main manager
		}

		if err := tx.Create(&teamManager).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to add manager: %v", err)
		}
	}

	// Add members
	for _, member := range req.Members {
		teamMember := model.Member{
			TeamID:     team.TeamID,
			MemberID:   member.MemberID,
			MemberName: member.MemberName,
		}

		if err := tx.Create(&teamMember).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to add member: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Load the complete team with managers and members
	return s.GetTeamByID(team.TeamID)
}

func (s *TeamService) GetTeamByID(teamID string) (*model.Team, error) {
	var team model.Team
	err := s.db.Preload("Managers").Preload("Members").First(&team, "team_id = ?", teamID).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *TeamService) AddMember(teamID string, req *model.AddMemberRequest, currentUserID string, token string) error {
	// Check if current user is a manager of this team
	if !s.isManagerOfTeam(currentUserID, teamID) {
		return errors.New("only managers can add members")
	}

	// Validate member exists in user service
	_, err := s.userServiceClient.ValidateUser(req.MemberID, token)
	if err != nil {
		return fmt.Errorf("member validation failed: %v", err)
	}

	// Check if user is already a member
	var existingMember model.Member
	result := s.db.Where("team_id = ? AND member_id = ?", teamID, req.MemberID).First(&existingMember)
	if result.Error == nil {
		return errors.New("user is already a member of this team")
	}

	// Check if user is already a manager
	var existingManager model.Manager
	result = s.db.Where("team_id = ? AND manager_id = ?", teamID, req.MemberID).First(&existingManager)
	if result.Error == nil {
		return errors.New("user is already a manager of this team")
	}

	// Add member
	member := model.Member{
		TeamID:     teamID,
		MemberID:   req.MemberID,
		MemberName: req.MemberName,
	}

	return s.db.Create(&member).Error
}

func (s *TeamService) RemoveMember(teamID, memberID, currentUserID string) error {
	// Check if current user is a manager of this team
	if !s.isManagerOfTeam(currentUserID, teamID) {
		return errors.New("only managers can remove members")
	}

	result := s.db.Where("team_id = ? AND member_id = ?", teamID, memberID).Delete(&model.Member{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found in team")
	}

	return nil
}

func (s *TeamService) AddManager(teamID string, req *model.AddManagerRequest, currentUserID string, token string) error {
	// Check if current user is the main manager of this team
	if !s.isMainManagerOfTeam(currentUserID, teamID) {
		return errors.New("only main manager can add other managers")
	}

	// Validate manager exists and has manager/admin role
	_, err := s.userServiceClient.ValidateRole(req.ManagerID, []string{"manager", "admin"}, token)
	if err != nil {
		return fmt.Errorf("manager validation failed: %v", err)
	}

	// Check if user is already a manager
	var existingManager model.Manager
	result := s.db.Where("team_id = ? AND manager_id = ?", teamID, req.ManagerID).First(&existingManager)
	if result.Error == nil {
		return errors.New("user is already a manager of this team")
	}

	// Remove from members if exists
	s.db.Where("team_id = ? AND member_id = ?", teamID, req.ManagerID).Delete(&model.Member{})

	// Add manager
	manager := model.Manager{
		TeamID:      teamID,
		ManagerID:   req.ManagerID,
		ManagerName: req.ManagerName,
		IsMain:      false, // Only one main manager per team
	}

	return s.db.Create(&manager).Error
}

func (s *TeamService) RemoveManager(teamID, managerID, currentUserID string) error {
	// Check if current user is the main manager of this team
	if !s.isMainManagerOfTeam(currentUserID, teamID) {
		return errors.New("only main manager can remove other managers")
	}

	// Cannot remove main manager
	var manager model.Manager
	result := s.db.Where("team_id = ? AND manager_id = ?", teamID, managerID).First(&manager)
	if result.Error != nil {
		return errors.New("manager not found in team")
	}

	if manager.IsMain {
		return errors.New("cannot remove main manager")
	}

	return s.db.Where("team_id = ? AND manager_id = ?", teamID, managerID).Delete(&model.Manager{}).Error
}

func (s *TeamService) isManagerOfTeam(userID, teamID string) bool {
	var count int64
	s.db.Model(&model.Manager{}).Where("team_id = ? AND manager_id = ?", teamID, userID).Count(&count)
	return count > 0
}

func (s *TeamService) isMainManagerOfTeam(userID, teamID string) bool {
	var count int64
	s.db.Model(&model.Manager{}).Where("team_id = ? AND manager_id = ? AND is_main = ?", teamID, userID, true).Count(&count)
	return count > 0
}

func (s *TeamService) GetAllTeams() ([]model.Team, error) {
	var teams []model.Team
	err := s.db.Preload("Managers").Preload("Members").Find(&teams).Error
	return teams, err
}

// Helper function to extract user ID from JWT token
func ExtractUserIDFromToken(tokenString string) (string, error) {
	// This is a simplified version - you should use proper JWT parsing
	// For now, assuming the token contains the user ID
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}
	
	// You should implement proper JWT parsing here
	// This is just a placeholder
	return "", errors.New("implement JWT parsing")
}