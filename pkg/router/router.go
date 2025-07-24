package router

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "golang/internal/team"
    "golang/internal/asset"
)

func Setup(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    teamGroup := r.Group("/teams")
    team.RegisterTeamRoutes(teamGroup, db)

    userGroup := r.Group("/assets")
    asset.RegisterRoutes(userGroup, db)


    return r
}
