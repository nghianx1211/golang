package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"asset-service/internal/model"
	"asset-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	userServiceClient *service.UserServiceClient
}

func NewAuthMiddleware(userServiceClient *service.UserServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		userServiceClient: userServiceClient,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token with user service
		userInfo, err := m.userServiceClient.ValidateToken(token)
		fmt.Println("user", userInfo)
		fmt.Println("err", err)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(userInfo.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", userID)
		c.Set("user_role", userInfo.Role)
		c.Set("user_info", userInfo)
		c.Set("token", token)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if userRole.(string) != "manager" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Manager access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper functions to get user info from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	return userID.(uuid.UUID), true
}

func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return userRole.(string), true
}

func GetUserInfo(c *gin.Context) (*model.UserInfo, bool) {
	userInfo, exists := c.Get("user_info")
	if !exists {
		return nil, false
	}
	return userInfo.(*model.UserInfo), true
}

func GetToken(c *gin.Context) (string, bool) {
	token, exists := c.Get("token")
	if !exists {
		return "", false
	}
	return token.(string), true
}