package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware: nếu có Authorization header → gắn claims vào context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := ParseAccessToken(tokenStr)
			if err == nil {
				// Token hợp lệ → gắn vào context
				ctx := WithUserID(c.Request.Context(), claims.UserID)
				ctx = WithRole(ctx, claims.Role)
				c.Request = c.Request.WithContext(ctx)
			}
			// Nếu token lỗi → bỏ qua, không chặn ở đây
		}

		c.Next()
	}
}
