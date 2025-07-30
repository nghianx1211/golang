package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	
	"github.com/google/uuid"
	"asset-service/internal/model"
	"github.com/joho/godotenv"

)

type UserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type FetchUsersResponse struct {
	FetchUsers []model.UserInfo `json:"fetchUsers"`
}

func NewUserServiceClient() *UserServiceClient {
	_ = godotenv.Load() 
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4000"
	}
	
	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *UserServiceClient) makeGraphQLRequest(query string, variables map[string]interface{}, token string) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
	
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", c.baseURL+"/query", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}
	
	return &gqlResp, nil
}

func (c *UserServiceClient) ValidateToken(token string) (*model.UserInfo, error) {
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
	
	resp, err := c.makeGraphQLRequest(query, nil, token)
	if err != nil {
		return nil, err
	}
	
	var fetchResp FetchUsersResponse
	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := json.Unmarshal(dataBytes, &fetchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	
	if len(fetchResp.FetchUsers) > 0 {
		return &fetchResp.FetchUsers[0], nil
	}
	
	return nil, fmt.Errorf("no user found")
}

func (c *UserServiceClient) GetUserInfo(userID uuid.UUID, token string) (*model.UserInfo, error) {
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
	
	resp, err := c.makeGraphQLRequest(query, nil, token)
	if err != nil {
		return nil, err
	}
	
	var fetchResp FetchUsersResponse
	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := json.Unmarshal(dataBytes, &fetchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	
	// Find the specific user
	for _, user := range fetchResp.FetchUsers {
		if user.UserID == userID.String() {
			return &user, nil
		}
	}
	
	return nil, fmt.Errorf("user not found")
}

func (c *UserServiceClient) CheckUserExists(userID uuid.UUID, token string) (bool, error) {
	_, err := c.GetUserInfo(userID, token)
	if err != nil {
		return false, nil // User doesn't exist or not accessible
	}
	return true, nil
}

func (c *UserServiceClient) GetTeamMembers(teamID uuid.UUID, token string) ([]model.UserInfo, error) {
	// This would need to be implemented based on your user service's team management
	// For now, returning all users as a placeholder
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
	
	resp, err := c.makeGraphQLRequest(query, nil, token)
	if err != nil {
		return nil, err
	}
	
	var fetchResp FetchUsersResponse
	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := json.Unmarshal(dataBytes, &fetchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	
	return fetchResp.FetchUsers, nil
}