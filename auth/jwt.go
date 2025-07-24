package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID uint
	Role   string
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
