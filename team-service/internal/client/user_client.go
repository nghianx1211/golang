package client

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type UserRoleResponse struct {
	Role string `json:"role"`
}

func CheckUserRole(userID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://user-service/users/%s/role", userID))
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user role")
	}
	defer resp.Body.Close()

	var result UserRoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Role, nil
}
