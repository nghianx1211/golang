

package asset

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

// Folder logic

func (s *Service) CreateFolder(name string, ownerID uint) (*model.Folder, error) {
    folder := &model.Folder{
        Name:    name,
        OwnerID: ownerID,
    }
    if err := s.DB.Create(folder).Error; err != nil {
        return nil, err
    }
    return folder, nil
}

func (s *Service) GetFolder(folderID uint) (*model.Folder, error) {
    var folder model.Folder
    if err := s.DB.Preload("Notes").First(&folder, folderID).Error; err != nil {
        return nil, err
    }
    return &folder, nil
}

func (s *Service) UpdateFolderName(folderID uint, name string) error {
    return s.DB.Model(&model.Folder{}).Where("folder_id = ?", folderID).Update("name", name).Error
}

func (s *Service) DeleteFolder(folderID uint) error {
    return s.DB.Delete(&model.Folder{}, folderID).Error
}

// Note logic

func (s *Service) CreateNote(folderID uint, title, body string, ownerID uint) (*model.Note, error) {
    note := &model.Note{
        Title:    title,
        Body:     body,
        FolderID: folderID,
        OwnerID:  ownerID,
    }
    if err := s.DB.Create(note).Error; err != nil {
        return nil, err
    }
    return note, nil
}

func (s *Service) GetNote(noteID uint) (*model.Note, error) {
    var note model.Note
    if err := s.DB.First(&note, noteID).Error; err != nil {
        return nil, err
    }
    return &note, nil
}

func (s *Service) UpdateNote(noteID uint, title, body string) error {
    return s.DB.Model(&model.Note{}).Where("note_id = ?", noteID).Updates(model.Note{Title: title, Body: body}).Error
}

func (s *Service) DeleteNote(noteID uint) error {
    return s.DB.Delete(&model.Note{}, noteID).Error
}

// Sharing

func (s *Service) ShareFolder(folderID, userID uint, access string) error {
    share := model.FolderShare{FolderID: folderID, UserID: userID, Access: access}
    return s.DB.Create(&share).Error
}

func (s *Service) RevokeFolderShare(folderID, userID uint) error {
    return s.DB.Delete(&model.FolderShare{}, "folder_id = ? AND user_id = ?", folderID, userID).Error
}

func (s *Service) ShareNote(noteID, userID uint, access string) error {
    share := model.NoteShare{NoteID: noteID, UserID: userID, Access: access}
    return s.DB.Create(&share).Error
}

func (s *Service) RevokeNoteShare(noteID, userID uint) error {
    return s.DB.Delete(&model.NoteShare{}, "note_id = ? AND user_id = ?", noteID, userID).Error
}

