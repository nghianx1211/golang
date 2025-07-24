package router

import (
	"github.com/nghianx1211/golang/internal/asset"
	"github.com/nghianx1211/golang/internal/team"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	teamGroup := r.Group("/teams")
	team.RegisterTeamRoutes(teamGroup, db)

	userGroup := r.Group("/assets")
	asset.RegisterRoutes(userGroup, db)

	return r
}
