package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserServiceClient struct {
	BaseURL string
	Client  *http.Client
}

type GraphQLRequest struct {
	Query string `json:"query"`
}

type UserData struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type FetchUsersResponse struct {
	Data struct {
		FetchUsers []UserData `json:"fetchUsers"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (u *UserServiceClient) FetchUsers(token string) ([]UserData, error) {
	query := `
		query {
			fetchUsers {
				userID
				username
				email
				role
			}
		}
	`

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", u.BaseURL+"/query", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// ✅ Nếu không phải 200 thì trả về toàn bộ body
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service error (status %d): %s", resp.StatusCode, string(body))
	}

	var response FetchUsersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v. raw body: %s", err, string(body))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("user service error: %s", response.Errors[0].Message)
	}

	return response.Data.FetchUsers, nil
}


func (u *UserServiceClient) ValidateUser(userID string, token string) (*UserData, error) {
	users, err := u.FetchUsers(token)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.UserID == userID {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (u *UserServiceClient) ValidateRole(userID string, expectedRoles []string, token string) (*UserData, error) {
	user, err := u.ValidateUser(userID, token)
	fmt.Println("user", user)
	fmt.Println("err", err)
	if err != nil {
		return nil, err
	}

	for _, role := range expectedRoles {
		if user.Role == role {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user does not have required role")
}