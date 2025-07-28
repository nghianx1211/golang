
package asset

import (
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
    service := NewService(db)
    handler := NewHandler(service)

    // Folder CRUD
    r.POST("/folders", handler.CreateFolder)
    r.GET("/folders/:folderId", handler.GetFolder)
    r.PUT("/folders/:folderId", handler.UpdateFolder)
    r.DELETE("/folders/:folderId", handler.DeleteFolder)

    // Note CRUD
    r.POST("/folders/:folderId/notes", handler.CreateNote)
    r.GET("/notes/:noteId", handler.GetNote)
    r.PUT("/notes/:noteId", handler.UpdateNote)
    r.DELETE("/notes/:noteId", handler.DeleteNote)

    // Sharing
    r.POST("/folders/:folderId/share", handler.ShareFolder)
    r.DELETE("/folders/:folderId/share/:userId", handler.RevokeFolderShare)
    r.POST("/notes/:noteId/share", handler.ShareNote)
    r.DELETE("/notes/:noteId/share/:userId", handler.RevokeNoteShare)
}
