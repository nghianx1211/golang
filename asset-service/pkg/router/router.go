
package router

import (
	"github.com/gin-gonic/gin"
	"asset-service/internal/handler"
	"asset-service/internal/middleware"
)

func SetupAssetRoutes(r *gin.Engine, assetHandler *handler.AssetHandler, authMiddleware *middleware.AuthMiddleware) {
	// Apply authentication middleware to all routes
	api := r.Group("/api/v1")
	api.Use(authMiddleware.RequireAuth())

	// Folder Management Routes
	folders := api.Group("/folders")
	{
		folders.POST("", assetHandler.CreateFolder)                          // POST /folders
		folders.GET("/:folderId", assetHandler.GetFolder)                    // GET /folders/:folderId
		folders.PUT("/:folderId", assetHandler.UpdateFolder)                 // PUT /folders/:folderId
		folders.DELETE("/:folderId", assetHandler.DeleteFolder)              // DELETE /folders/:folderId
		
		// Notes within folders
		folders.POST("/:folderId/notes", assetHandler.CreateNote)            // POST /folders/:folderId/notes
		
		// Folder sharing
		folders.POST("/:folderId/share", assetHandler.ShareFolder)           // POST /folders/:folderId/share
		folders.DELETE("/:folderId/share/:userId", assetHandler.RevokeFolderSharing) // DELETE /folders/:folderId/share/:userId
	}

	// Note Management Routes
	notes := api.Group("/notes")
	{
		notes.GET("/:noteId", assetHandler.GetNote)                          // GET /notes/:noteId
		notes.PUT("/:noteId", assetHandler.UpdateNote)                       // PUT /notes/:noteId
		notes.DELETE("/:noteId", assetHandler.DeleteNote)                    // DELETE /notes/:noteId
		
		// Note sharing
		notes.POST("/:noteId/share", assetHandler.ShareNote)                 // POST /notes/:noteId/share
		notes.DELETE("/:noteId/share/:userId", assetHandler.RevokeNoteSharing) // DELETE /notes/:noteId/share/:userId
	}

	// Manager-only Routes
	manager := api.Group("")
	manager.Use(authMiddleware.RequireManager())
	{
		manager.GET("/teams/:teamId/assets", assetHandler.GetTeamAssets)     // GET /teams/:teamId/assets
		manager.GET("/users/:userId/assets", assetHandler.GetUserAssets)     // GET /users/:userId/assets
	}
}