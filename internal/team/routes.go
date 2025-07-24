package team

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterTeamRoutes(r *gin.RouterGroup, db *gorm.DB) {
	service := NewService(db)
	handler := NewHandler(service)

	r.POST("", handler.CreateTeam)
	r.POST("/:teamId/members", handler.AddMember)
	r.DELETE("/:teamId/members/:memberId", handler.RemoveMember)
	r.POST("/:teamId/managers", handler.AddManager)
	r.DELETE("/:teamId/managers/:managerId", handler.RemoveManager)
}