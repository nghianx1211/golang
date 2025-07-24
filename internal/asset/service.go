package asset

import (
	"github.com/nghianx1211/golang/internal/model"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// ---------- Folder logic ----------

func (s *Service) CreateFolder(name string, ownerID string) (*model.Folder, error) {
	folder := &model.Folder{
		Name:    name,
		OwnerID: ownerID,
	}
	if err := s.DB.Create(folder).Error; err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *Service) GetFolder(folderID string) (*model.Folder, error) {
	var folder model.Folder
	if err := s.DB.Preload("Notes").First(&folder, "folder_id = ?", folderID).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (s *Service) UpdateFolderName(folderID string, name string) error {
	return s.DB.Model(&model.Folder{}).Where("folder_id = ?", folderID).Update("name", name).Error
}

func (s *Service) DeleteFolder(folderID string) error {
	return s.DB.Delete(&model.Folder{}, "folder_id = ?", folderID).Error
}

// ---------- Note logic ----------

func (s *Service) CreateNote(folderID string, title, body string, ownerID string) (*model.Note, error) {
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

func (s *Service) GetNote(noteID string) (*model.Note, error) {
	var note model.Note
	if err := s.DB.First(&note, "note_id = ?", noteID).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *Service) UpdateNote(noteID string, title, body string) error {
	return s.DB.Model(&model.Note{}).
		Where("note_id = ?", noteID).
		Updates(map[string]interface{}{
			"title": title,
			"body":  body,
		}).Error
}

func (s *Service) DeleteNote(noteID string) error {
	return s.DB.Delete(&model.Note{}, "note_id = ?", noteID).Error
}

// ---------- Sharing logic ----------

func (s *Service) ShareFolder(folderID, userID string, access string) error {
	share := model.FolderShare{
		FolderID: folderID,
		UserID:   userID,
		Access:   access,
	}
	return s.DB.Create(&share).Error
}

func (s *Service) RevokeFolderShare(folderID, userID string) error {
	return s.DB.Delete(&model.FolderShare{}, "folder_id = ? AND user_id = ?", folderID, userID).Error
}

func (s *Service) ShareNote(noteID, userID string, access string) error {
	share := model.NoteShare{
		NoteID: noteID,
		UserID: userID,
		Access: access,
	}
	return s.DB.Create(&share).Error
}

func (s *Service) RevokeNoteShare(noteID, userID string) error {
	return s.DB.Delete(&model.NoteShare{}, "note_id = ? AND user_id = ?", noteID, userID).Error
}
