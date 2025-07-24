package user

import (
	"gorm.io/gorm"
	"golang/internal/model"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) GetUserAssets(userID uint) (folders []model.Folder, folderShares []model.FolderShare, noteShares []model.NoteShare, err error) {
    if err = s.DB.Where("owner_id = ?", userID).Preload("Notes").Find(&folders).Error; err != nil {
        return
    }

    if err = s.DB.Where("user_id = ?", userID).Find(&folderShares).Error; err != nil {
        return
    }

    if err = s.DB.Where("user_id = ?", userID).Find(&noteShares).Error; err != nil {
        return
    }

    return
}