package service

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"asset-service/internal/model"
)

type AssetService struct {
	db               *gorm.DB
	userServiceClient *UserServiceClient
}

func NewAssetService(db *gorm.DB, userServiceClient *UserServiceClient) *AssetService {
	return &AssetService{
		db:               db,
		userServiceClient: userServiceClient,
	}
}

// Folder CRUD Operations
func (s *AssetService) CreateFolder(req *model.CreateFolderRequest, ownerID uuid.UUID) (*model.FolderResponse, error) {
	folder := &model.Folder{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	}

	if err := s.db.Create(folder).Error; err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return s.folderToResponse(folder), nil
}

func (s *AssetService) GetFolder(folderID uuid.UUID, userID uuid.UUID, userRole string) (*model.FolderResponse, error) {
	var folder model.Folder
	
	// Check if user can access this folder
	query := s.db.Where("id = ?", folderID)
	
	// If not owner and not manager, check if folder is shared with user
	if userRole != "manager" {
		query = query.Where("owner_id = ? OR id IN (SELECT folder_id FROM folder_shares WHERE user_id = ?)", userID, userID)
	}
	
	if err := query.Preload("Notes").Preload("SharedWith").First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("folder not found or access denied")
		}
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return s.folderToResponse(&folder), nil
}

func (s *AssetService) UpdateFolder(folderID uuid.UUID, req *model.UpdateFolderRequest, userID uuid.UUID, userRole string) (*model.FolderResponse, error) {
	var folder model.Folder
	
	// Check permissions - only owner or users with write access can update
	query := s.db.Where("id = ?", folderID)
	if userRole != "manager" {
		query = query.Where("owner_id = ? OR id IN (SELECT folder_id FROM folder_shares WHERE user_id = ? AND permission = 'write')", userID, userID)
	}
	
	if err := query.First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("folder not found or access denied")
		}
		return nil, fmt.Errorf("failed to find folder: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		folder.Name = req.Name
	}
	if req.Description != "" {
		folder.Description = req.Description
	}

	if err := s.db.Save(&folder).Error; err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return s.folderToResponse(&folder), nil
}

func (s *AssetService) DeleteFolder(folderID uuid.UUID, userID uuid.UUID, userRole string) error {
	var folder model.Folder
	
	// Only owner can delete folder
	query := s.db.Where("id = ? AND owner_id = ?", folderID, userID)
	
	if err := query.First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("folder not found or access denied")
		}
		return fmt.Errorf("failed to find folder: %w", err)
	}

	// Delete folder (cascade will handle notes and shares)
	if err := s.db.Delete(&folder).Error; err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}

// Note CRUD Operations
func (s *AssetService) CreateNote(folderID uuid.UUID, req *model.CreateNoteRequest, userID uuid.UUID, userRole string) (*model.NoteResponse, error) {
	// Check if user can create notes in this folder
	var folder model.Folder
	query := s.db.Where("id = ?", folderID)
	
	if userRole != "manager" {
		query = query.Where("owner_id = ? OR id IN (SELECT folder_id FROM folder_shares WHERE user_id = ? AND permission = 'write')", userID, userID)
	}
	
	if err := query.First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("folder not found or access denied")
		}
		return nil, fmt.Errorf("failed to find folder: %w", err)
	}

	note := &model.Note{
		Title:    req.Title,
		Content:  req.Content,
		FolderID: folderID,
		OwnerID:  userID,
	}

	if err := s.db.Create(note).Error; err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return s.noteToResponse(note), nil
}

func (s *AssetService) GetNote(noteID uuid.UUID, userID uuid.UUID, userRole string) (*model.NoteResponse, error) {
	var note model.Note
	
	query := s.db.Where("id = ?", noteID)
	
	// If not manager, check permissions
	if userRole != "manager" {
		query = query.Where(`owner_id = ? OR 
			folder_id IN (SELECT folder_id FROM folder_shares WHERE user_id = ?) OR
			id IN (SELECT note_id FROM note_shares WHERE user_id = ?)`, userID, userID, userID)
	}
	
	if err := query.Preload("SharedWith").First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("note not found or access denied")
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return s.noteToResponse(&note), nil
}

func (s *AssetService) UpdateNote(noteID uuid.UUID, req *model.UpdateNoteRequest, userID uuid.UUID, userRole string) (*model.NoteResponse, error) {
	var note model.Note
	
	query := s.db.Where("id = ?", noteID)
	
	// Check write permissions
	if userRole != "manager" {
		query = query.Where(`owner_id = ? OR 
			folder_id IN (SELECT folder_id FROM folder_shares WHERE user_id = ? AND permission = 'write') OR
			id IN (SELECT note_id FROM note_shares WHERE user_id = ? AND permission = 'write')`, userID, userID, userID)
	}
	
	if err := query.First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("note not found or access denied")
		}
		return nil, fmt.Errorf("failed to find note: %w", err)
	}

	// Update fields if provided
	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}

	if err := s.db.Save(&note).Error; err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return s.noteToResponse(&note), nil
}

func (s *AssetService) DeleteNote(noteID uuid.UUID, userID uuid.UUID, userRole string) error {
	var note model.Note
	
	// Only owner can delete note
	query := s.db.Where("id = ? AND owner_id = ?", noteID, userID)
	
	if err := query.First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("note not found or access denied")
		}
		return fmt.Errorf("failed to find note: %w", err)
	}

	if err := s.db.Delete(&note).Error; err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}

// Sharing Operations
func (s *AssetService) ShareFolder(folderID uuid.UUID, req *model.ShareRequest, sharedBy uuid.UUID, token string) error {
	// Check if folder exists and user is owner
	var folder model.Folder
	if err := s.db.Where("id = ? AND owner_id = ?", folderID, sharedBy).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("folder not found or access denied")
		}
		return fmt.Errorf("failed to find folder: %w", err)
	}

	// Validate target user exists
	exists, err := s.userServiceClient.CheckUserExists(req.UserID, token)
	if err != nil || !exists {
		return fmt.Errorf("target user not found")
	}

	// Check if already shared
	var existingShare model.FolderShare
	if err := s.db.Where("folder_id = ? AND user_id = ?", folderID, req.UserID).First(&existingShare).Error; err == nil {
		// Update existing permission
		existingShare.Permission = req.Permission
		return s.db.Save(&existingShare).Error
	}

	// Create new share
	share := &model.FolderShare{
		FolderID:   folderID,
		UserID:     req.UserID,
		Permission: req.Permission,
		SharedBy:   sharedBy,
	}

	return s.db.Create(share).Error
}

func (s *AssetService) RevokeFolderSharing(folderID uuid.UUID, targetUserID uuid.UUID, ownerID uuid.UUID) error {
	// Check if folder exists and user is owner
	var folder model.Folder
	if err := s.db.Where("id = ? AND owner_id = ?", folderID, ownerID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("folder not found or access denied")
		}
		return fmt.Errorf("failed to find folder: %w", err)
	}

	return s.db.Where("folder_id = ? AND user_id = ?", folderID, targetUserID).Delete(&model.FolderShare{}).Error
}

func (s *AssetService) ShareNote(noteID uuid.UUID, req *model.ShareRequest, sharedBy uuid.UUID, token string) error {
	// Check if note exists and user is owner
	var note model.Note
	if err := s.db.Where("id = ? AND owner_id = ?", noteID, sharedBy).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("note not found or access denied")
		}
		return fmt.Errorf("failed to find note: %w", err)
	}

	// Validate target user exists
	exists, err := s.userServiceClient.CheckUserExists(req.UserID, token)
	if err != nil || !exists {
		return fmt.Errorf("target user not found")
	}

	// Check if already shared
	var existingShare model.NoteShare
	if err := s.db.Where("note_id = ? AND user_id = ?", noteID, req.UserID).First(&existingShare).Error; err == nil {
		// Update existing permission
		existingShare.Permission = req.Permission
		return s.db.Save(&existingShare).Error
	}

	// Create new share
	share := &model.NoteShare{
		NoteID:     noteID,
		UserID:     req.UserID,
		Permission: req.Permission,
		SharedBy:   sharedBy,
	}

	return s.db.Create(share).Error
}

func (s *AssetService) RevokeNoteSharing(noteID uuid.UUID, targetUserID uuid.UUID, ownerID uuid.UUID) error {
	// Check if note exists and user is owner
	var note model.Note
	if err := s.db.Where("id = ? AND owner_id = ?", noteID, ownerID).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("note not found or access denied")
		}
		return fmt.Errorf("failed to find note: %w", err)
	}

	return s.db.Where("note_id = ? AND user_id = ?", noteID, targetUserID).Delete(&model.NoteShare{}).Error
}

// Manager-only Operations
func (s *AssetService) GetTeamAssets(teamID uuid.UUID, token string) (*model.AssetResponse, error) {
	// Get team members
	members, err := s.userServiceClient.GetTeamMembers(teamID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	var memberIDs []uuid.UUID
	for _, member := range members {
		if id, err := uuid.Parse(member.UserID); err == nil {
			memberIDs = append(memberIDs, id)
		}
	}

	// Get all folders owned by or shared with team members
	var folders []model.Folder
	s.db.Where("owner_id IN ? OR id IN (SELECT folder_id FROM folder_shares WHERE user_id IN ?)", memberIDs, memberIDs).
		Preload("Notes").Preload("SharedWith").Find(&folders)

	// Get all notes owned by or shared with team members
	var notes []model.Note
	s.db.Where("owner_id IN ? OR id IN (SELECT note_id FROM note_shares WHERE user_id IN ?)", memberIDs, memberIDs).
		Preload("SharedWith").Find(&notes)

	folderResponses := make([]model.FolderResponse, len(folders))
	for i, folder := range folders {
		folderResponses[i] = *s.folderToResponse(&folder)
	}

	noteResponses := make([]model.NoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = *s.noteToResponse(&note)
	}

	return &model.AssetResponse{
		Folders: folderResponses,
		Notes:   noteResponses,
	}, nil
}

func (s *AssetService) GetUserAssets(targetUserID uuid.UUID, requestorRole string) (*model.AssetResponse, error) {
	if requestorRole != "manager" {
		return nil, fmt.Errorf("access denied")
	}

	// Get folders owned by or shared with user
	var folders []model.Folder
	s.db.Where("owner_id = ? OR id IN (SELECT folder_id FROM folder_shares WHERE user_id = ?)", targetUserID, targetUserID).
		Preload("Notes").Preload("SharedWith").Find(&folders)

	// Get notes owned by or shared with user
	var notes []model.Note
	s.db.Where("owner_id = ? OR id IN (SELECT note_id FROM note_shares WHERE user_id = ?)", targetUserID, targetUserID).
		Preload("SharedWith").Find(&notes)

	folderResponses := make([]model.FolderResponse, len(folders))
	for i, folder := range folders {
		folderResponses[i] = *s.folderToResponse(&folder)
	}

	noteResponses := make([]model.NoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = *s.noteToResponse(&note)
	}

	return &model.AssetResponse{
		Folders: folderResponses,
		Notes:   noteResponses,
	}, nil
}

// Helper methods
func (s *AssetService) folderToResponse(folder *model.Folder) *model.FolderResponse {
	resp := &model.FolderResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		OwnerID:     folder.OwnerID,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}

	// Convert notes
	if folder.Notes != nil {
		resp.Notes = make([]model.NoteResponse, len(folder.Notes))
		for i, note := range folder.Notes {
			resp.Notes[i] = *s.noteToResponse(&note)
		}
	}

	// Convert shares
	if folder.SharedWith != nil {
		resp.SharedWith = make([]model.ShareInfo, len(folder.SharedWith))
		for i, share := range folder.SharedWith {
			resp.SharedWith[i] = model.ShareInfo{
				UserID:     share.UserID,
				Permission: share.Permission,
				SharedBy:   share.SharedBy,
				CreatedAt:  share.CreatedAt,
			}
		}
	}

	return resp
}

func (s *AssetService) noteToResponse(note *model.Note) *model.NoteResponse {
	resp := &model.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		FolderID:  note.FolderID,
		OwnerID:   note.OwnerID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	// Convert shares
	if note.SharedWith != nil {
		resp.SharedWith = make([]model.ShareInfo, len(note.SharedWith))
		for i, share := range note.SharedWith {
			resp.SharedWith[i] = model.ShareInfo{
				UserID:     share.UserID,
				Permission: share.Permission,
				SharedBy:   share.SharedBy,
				CreatedAt:  share.CreatedAt,
			}
		}
	}

	return resp
}