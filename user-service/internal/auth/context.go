package auth

import (
	"context"
	"errors"
	"fmt"
)

type contextKey string

const (
	userIDKey = contextKey("userID")
	roleKey   = contextKey("role")
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", errors.New("unauthenticated")
	}
	return userID, nil
}

func GetRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(roleKey).(string)
	fmt.Println("rolee", role)
	if !ok {
		return "", errors.New("unauthenticated")
	}
	return role, nil
}