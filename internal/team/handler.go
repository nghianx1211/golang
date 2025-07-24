package team

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

func (h *Handler) CreateTeam(c *gin.Context) {
	var input struct {
		TeamName     string `json:"teamName" binding:"required"`
		MainManagerID uint   `json:"mainManagerId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.Service.CreateTeam(input.TeamName, input.MainManagerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *Handler) AddMember(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	var input struct {
		UserID uint `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.AddMember(uint(teamID), input.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}

func (h *Handler) RemoveMember(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	memberID, _ := strconv.Atoi(c.Param("memberId"))

	if err := h.Service.RemoveMember(uint(teamID), uint(memberID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

func (h *Handler) AddManager(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	var input struct {
		RequesterID uint `json:"requesterId" binding:"required"`
		UserID      uint `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.AddManager(uint(teamID), input.RequesterID, input.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager added"})
}

func (h *Handler) RemoveManager(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	managerID, _ := strconv.Atoi(c.Param("managerId"))

	if err := h.Service.RemoveManager(uint(teamID), uint(managerID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed"})
}

func (h *Handler) GetTeamAssets(c *gin.Context) {
    teamID, _ := strconv.Atoi(c.Param("teamId"))
    assets, err := h.Service.GetTeamAssets(uint(teamID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, assets)
}
