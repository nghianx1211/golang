package team

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

func (h *Handler) CreateTeam(c *gin.Context) {
	var input struct {
		TeamName      string `json:"teamName" binding:"required"`
		MainManagerID string `json:"mainManagerId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := uuid.Parse(input.MainManagerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID for mainManagerId"})
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
	teamID := c.Param("teamId")
	var input struct {
		UserID string `json:"userId" binding:"required"`
	}

	if _, err := uuid.Parse(teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamId UUID"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := uuid.Parse(input.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId UUID"})
		return
	}

	if err := h.Service.AddMember(teamID, input.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}

func (h *Handler) RemoveMember(c *gin.Context) {
	teamID := c.Param("teamId")
	memberID := c.Param("memberId")

	if _, err := uuid.Parse(teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamId UUID"})
		return
	}
	if _, err := uuid.Parse(memberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid memberId UUID"})
		return
	}

	if err := h.Service.RemoveMember(teamID, memberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

func (h *Handler) AddManager(c *gin.Context) {
	teamID := c.Param("teamId")
	var input struct {
		RequesterID string `json:"requesterId" binding:"required"`
		UserID      string `json:"userId" binding:"required"`
	}

	if _, err := uuid.Parse(teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamId UUID"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := uuid.Parse(input.RequesterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requesterId UUID"})
		return
	}
	if _, err := uuid.Parse(input.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId UUID"})
		return
	}

	if err := h.Service.AddManager(teamID, input.RequesterID, input.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager added"})
}

func (h *Handler) RemoveManager(c *gin.Context) {
	teamID := c.Param("teamId")
	managerID := c.Param("managerId")

	if _, err := uuid.Parse(teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamId UUID"})
		return
	}
	if _, err := uuid.Parse(managerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid managerId UUID"})
		return
	}

	if err := h.Service.RemoveManager(teamID, managerID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed"})
}

func (h *Handler) GetTeamAssets(c *gin.Context) {
	teamID := c.Param("teamId")

	if _, err := uuid.Parse(teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamId UUID"})
		return
	}

	assets, err := h.Service.GetTeamAssets(teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}
