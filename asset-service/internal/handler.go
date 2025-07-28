package asset

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type Handler struct {
    Service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{Service: service}
}

// ---------- Folder ----------

func (h *Handler) CreateFolder(c *gin.Context) {
    var input struct {
        Name    string `json:"name" binding:"required"`
        OwnerID string `json:"ownerId" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if _, err := uuid.Parse(input.OwnerID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ownerId UUID"})
        return
    }

    folder, err := h.Service.CreateFolder(input.Name, input.OwnerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, folder)
}

func (h *Handler) GetFolder(c *gin.Context) {
    folderID := c.Param("folderId")
    if _, err := uuid.Parse(folderID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folderId UUID"})
        return
    }

    folder, err := h.Service.GetFolder(folderID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, folder)
}

func (h *Handler) UpdateFolder(c *gin.Context) {
    folderID := c.Param("folderId")
    if _, err := uuid.Parse(folderID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folderId UUID"})
        return
    }

    var input struct {
        Name string `json:"name"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Service.UpdateFolderName(folderID, input.Name); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Folder updated"})
}

func (h *Handler) DeleteFolder(c *gin.Context) {
    folderID := c.Param("folderId")
    if _, err := uuid.Parse(folderID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folderId UUID"})
        return
    }

    h.Service.DeleteFolder(folderID)
    c.JSON(http.StatusOK, gin.H{"message": "Folder deleted"})
}

// ---------- Note ----------

func (h *Handler) CreateNote(c *gin.Context) {
    folderID := c.Param("folderId")
    if _, err := uuid.Parse(folderID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folderId UUID"})
        return
    }

    var input struct {
        Title   string `json:"title"`
        Body    string `json:"body"`
        OwnerID string `json:"ownerId"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if _, err := uuid.Parse(input.OwnerID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ownerId UUID"})
        return
    }

    note, err := h.Service.CreateNote(folderID, input.Title, input.Body, input.OwnerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, note)
}

func (h *Handler) GetNote(c *gin.Context) {
    noteID := c.Param("noteId")
    if _, err := uuid.Parse(noteID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid noteId UUID"})
        return
    }

    note, err := h.Service.GetNote(noteID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (h *Handler) UpdateNote(c *gin.Context) {
    noteID := c.Param("noteId")
    if _, err := uuid.Parse(noteID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid noteId UUID"})
        return
    }

    var input struct {
        Title string `json:"title"`
        Body  string `json:"body"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.Service.UpdateNote(noteID, input.Title, input.Body)
    c.JSON(http.StatusOK, gin.H{"message": "Note updated"})
}

func (h *Handler) DeleteNote(c *gin.Context) {
    noteID := c.Param("noteId")
    if _, err := uuid.Parse(noteID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid noteId UUID"})
        return
    }

    h.Service.DeleteNote(noteID)
    c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

// ---------- Sharing ----------

func (h *Handler) ShareFolder(c *gin.Context) {
    folderID := c.Param("folderId")
    if _, err := uuid.Parse(folderID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folderId UUID"})
        return
    }

    var input struct {
        UserID string `json:"userId"`
        Access string `json:"access"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if _, err := uuid.Parse(input.UserID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId UUID"})
        return
    }

    h.Service.ShareFolder(folderID, input.UserID, input.Access)
    c.JSON(http.StatusOK, gin.H{"message": "Folder shared"})
}

func (h *Handler) RevokeFolderShare(c *gin.Context) {
    folderID := c.Param("folderId")
    userID := c.Param("userId")

    if _, err := uuid.Parse(folderID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folderId UUID"})
        return
    }
    if _, err := uuid.Parse(userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId UUID"})
        return
    }

    h.Service.RevokeFolderShare(folderID, userID)
    c.JSON(http.StatusOK, gin.H{"message": "Folder share revoked"})
}

func (h *Handler) ShareNote(c *gin.Context) {
    noteID := c.Param("noteId")
    if _, err := uuid.Parse(noteID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid noteId UUID"})
        return
    }

    var input struct {
        UserID string `json:"userId"`
        Access string `json:"access"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if _, err := uuid.Parse(input.UserID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId UUID"})
        return
    }

    h.Service.ShareNote(noteID, input.UserID, input.Access)
    c.JSON(http.StatusOK, gin.H{"message": "Note shared"})
}

func (h *Handler) RevokeNoteShare(c *gin.Context) {
    noteID := c.Param("noteId")
    userID := c.Param("userId")

    if _, err := uuid.Parse(noteID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid noteId UUID"})
        return
    }
    if _, err := uuid.Parse(userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId UUID"})
        return
    }

    h.Service.RevokeNoteShare(noteID, userID)
    c.JSON(http.StatusOK, gin.H{"message": "Note share revoked"})
}
