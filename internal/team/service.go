package team

import (
	"errors"
	"gorm.io/gorm"
	"golang/internal/model"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateTeam(teamName string, mainManagerID uint) (*model.Team, error) {
	var manager model.User
	if err := s.DB.First(&manager, mainManagerID).Error; err != nil {
		return nil, errors.New("main manager not found")
	}

	team := &model.Team{
		TeamName: teamName,
		Managers: []model.User{manager},
	}

	if err := s.DB.Create(team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (s *Service) AddMember(teamID, userID uint) error {
	var team model.Team
	if err := s.DB.Preload("Members").First(&team, teamID).Error; err != nil {
		return errors.New("team not found")
	}

	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	return s.DB.Model(&team).Association("Members").Append(&user)
}

func (s *Service) RemoveMember(teamID, userID uint) error {
	var team model.Team
	if err := s.DB.Preload("Members").First(&team, teamID).Error; err != nil {
		return errors.New("team not found")
	}

	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	return s.DB.Model(&team).Association("Members").Delete(&user)
}

func (s *Service) AddManager(teamID, requesterID, userID uint) error {
	var team model.Team
	if err := s.DB.Preload("Managers").First(&team, teamID).Error; err != nil {
		return errors.New("team not found")
	}

	if len(team.Managers) == 0 || team.Managers[0].UserID != requesterID {
		return errors.New("only main manager can add managers")
	}

	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	return s.DB.Model(&team).Association("Managers").Append(&user)
}

func (s *Service) RemoveManager(teamID, managerID uint) error {
	var team model.Team
	if err := s.DB.Preload("Managers").First(&team, teamID).Error; err != nil {
		return errors.New("team not found")
	}

	if len(team.Managers) == 0 {
		return errors.New("no managers found")
	}

	if team.Managers[0].UserID == managerID {
		return errors.New("cannot remove main manager")
	}

	var user model.User
	if err := s.DB.First(&user, managerID).Error; err != nil {
		return errors.New("manager not found")
	}

	return s.DB.Model(&team).Association("Managers").Delete(&user)
}

func (s *Service) GetTeamAssets(teamID uint) ([]model.Folder, error) {
    var team model.Team
    if err := s.DB.Preload("Members").First(&team, teamID).Error; err != nil {
        return nil, err
    }

    var memberIDs []uint
    for _, m := range team.Members {
        memberIDs = append(memberIDs, m.UserID)
    }

    var assets []model.Folder
    if err := s.DB.Where("owner_id IN ?", memberIDs).Preload("Notes").Find(&assets).Error; err != nil {
        return nil, err
    }

    return assets, nil
}
