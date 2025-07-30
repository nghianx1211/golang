package router

import (
	"team-service/internal/handler"   
	"team-service/internal/middleware" 

	"github.com/gin-gonic/gin"
)

// SetupTeamRoutes sets up all team-related routes
func SetupTeamRouter(router *gin.Engine, teamHandler *handler.TeamHandler, jwtSecret string) {
	// Apply auth middleware to all team routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtSecret))

	// Team routes
	teams := api.Group("/teams")
	{
		// Create team - requires manager or admin role
		teams.POST("", middleware.RequireRole("manager", "admin"), teamHandler.CreateTeam)
		
		// Get all teams - any authenticated user
		teams.GET("", teamHandler.GetAllTeams)
		
		// Get specific team - any authenticated user
		teams.GET("/:teamId", teamHandler.GetTeam)
		
		// Member management routes
		teams.POST("/:teamId/members", teamHandler.AddMember)
		teams.DELETE("/:teamId/members/:memberId", teamHandler.RemoveMember)
		
		// Manager management routes
		teams.POST("/:teamId/managers", teamHandler.AddManager)
		teams.DELETE("/:teamId/managers/:managerId", teamHandler.RemoveManager)
	}
}

// Alternative setup for direct routing (without middleware groups)
func SetupTeamRoutesSimple(router *gin.Engine, teamHandler *handler.TeamHandler) {
	// Team routes without auth middleware (if you handle auth differently)
	router.POST("/teams", teamHandler.CreateTeam)
	router.GET("/teams", teamHandler.GetAllTeams)
	router.GET("/teams/:teamId", teamHandler.GetTeam)
	router.POST("/teams/:teamId/members", teamHandler.AddMember)
	router.DELETE("/teams/:teamId/members/:memberId", teamHandler.RemoveMember)
	router.POST("/teams/:teamId/managers", teamHandler.AddManager)
	router.DELETE("/teams/:teamId/managers/:managerId", teamHandler.RemoveManager)
}