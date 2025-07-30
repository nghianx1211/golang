package handler

import (
	"net/http"
	"strings"

	"team-service/internal/model"   // Replace with your actual module path
	"team-service/internal/service" // Replace with your actual module path

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// CreateTeam handles POST /teams
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req model.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}
	token := authHeader[7:] // Remove "Bearer " prefix

	team, err := h.teamService.CreateTeam(&req, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// GetTeam handles GET /teams/:teamId
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamID := c.Param("teamId")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is required"})
		return
	}

	team, err := h.teamService.GetTeamByID(teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetAllTeams handles GET /teams
func (h *TeamHandler) GetAllTeams(c *gin.Context) {
	teams, err := h.teamService.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// AddMember handles POST /teams/:teamId/members
func (h *TeamHandler) AddMember(c *gin.Context) {
	teamID := c.Param("teamId")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is required"})
		return
	}

	var req model.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract token and current user ID
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}
	token := authHeader[7:]

	// Get current user ID from context (should be set by auth middleware)
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	err := h.teamService.AddMember(teamID, &req, currentUserID.(string), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// RemoveMember handles DELETE /teams/:teamId/members/:memberId
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID := c.Param("teamId")
	memberID := c.Param("memberId")

	if teamID == "" || memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId and memberId are required"})
		return
	}

	// Get current user ID from context
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	err := h.teamService.RemoveMember(teamID, memberID, currentUserID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// AddManager handles POST /teams/:teamId/managers
func (h *TeamHandler) AddManager(c *gin.Context) {
	teamID := c.Param("teamId")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is required"})
		return
	}

	var req model.AddManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}
	token := authHeader[7:]

	// Get current user ID from context
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	err := h.teamService.AddManager(teamID, &req, currentUserID.(string), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager added successfully"})
}

// RemoveManager handles DELETE /teams/:teamId/managers/:managerId
func (h *TeamHandler) RemoveManager(c *gin.Context) {
	teamID := c.Param("teamId")
	managerID := c.Param("managerId")

	if teamID == "" || managerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId and managerId are required"})
		return
	}

	// Get current user ID from context
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	err := h.teamService.RemoveManager(teamID, managerID, currentUserID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed successfully"})
}