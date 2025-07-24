package asset

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    Service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{Service: service}
}

// Folder

func (h *Handler) CreateFolder(c *gin.Context) {
    var input struct {
        Name    string `json:"name" binding:"required"`
        OwnerID uint   `json:"ownerId" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
    id, _ := strconv.Atoi(c.Param("folderId"))
    folder, err := h.Service.GetFolder(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, folder)
}

func (h *Handler) UpdateFolder(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("folderId"))
    var input struct{ Name string `json:"name"` }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.Service.UpdateFolderName(uint(id), input.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Folder updated"})
}

func (h *Handler) DeleteFolder(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("folderId"))
    h.Service.DeleteFolder(uint(id))
    c.JSON(http.StatusOK, gin.H{"message": "Folder deleted"})
}

// Note

func (h *Handler) CreateNote(c *gin.Context) {
    folderID, _ := strconv.Atoi(c.Param("folderId"))
    var input struct {
        Title   string `json:"title"`
        Body    string `json:"body"`
        OwnerID uint   `json:"ownerId"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    note, err := h.Service.CreateNote(uint(folderID), input.Title, input.Body, input.OwnerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, note)
}

func (h *Handler) GetNote(c *gin.Context) {
    noteID, _ := strconv.Atoi(c.Param("noteId"))
    note, err := h.Service.GetNote(uint(noteID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, note)
}

func (h *Handler) UpdateNote(c *gin.Context) {
    noteID, _ := strconv.Atoi(c.Param("noteId"))
    var input struct{ Title, Body string }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.Service.UpdateNote(uint(noteID), input.Title, input.Body)
    c.JSON(http.StatusOK, gin.H{"message": "Note updated"})
}

func (h *Handler) DeleteNote(c *gin.Context) {
    noteID, _ := strconv.Atoi(c.Param("noteId"))
    h.Service.DeleteNote(uint(noteID))
    c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

// Sharing

func (h *Handler) ShareFolder(c *gin.Context) {
    folderID, _ := strconv.Atoi(c.Param("folderId"))
    var input struct {
        UserID uint   `json:"userId"`
        Access string `json:"access"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.Service.ShareFolder(uint(folderID), input.UserID, input.Access)
    c.JSON(http.StatusOK, gin.H{"message": "Folder shared"})
}

func (h *Handler) RevokeFolderShare(c *gin.Context) {
    folderID, _ := strconv.Atoi(c.Param("folderId"))
    userID, _ := strconv.Atoi(c.Param("userId"))

    h.Service.RevokeFolderShare(uint(folderID), uint(userID))
    c.JSON(http.StatusOK, gin.H{"message": "Folder share revoked"})
}

func (h *Handler) ShareNote(c *gin.Context) {
    noteID, _ := strconv.Atoi(c.Param("noteId"))
    var input struct {
        UserID uint   `json:"userId"`
        Access string `json:"access"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.Service.ShareNote(uint(noteID), input.UserID, input.Access)
    c.JSON(http.StatusOK, gin.H{"message": "Note shared"})
}

func (h *Handler) RevokeNoteShare(c *gin.Context) {
    noteID, _ := strconv.Atoi(c.Param("noteId"))
    userID, _ := strconv.Atoi(c.Param("userId"))

    h.Service.RevokeNoteShare(uint(noteID), uint(userID))
    c.JSON(http.StatusOK, gin.H{"message": "Note share revoked"})
}