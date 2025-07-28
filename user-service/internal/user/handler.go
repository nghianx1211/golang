package user

import (
    "context"
)

type Resolver struct {
    Service *Service
}

func (r *Resolver) CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
    user, err := r.Service.CreateUser(input)
    if err != nil {
        return nil, err
    }

    return &UserOutput{
        UserID:   user.UserID,
        Username: user.Username,
        Email:    user.Email,
        Role:     user.Role,
    }, nil
}

func (r *Resolver) Login(ctx context.Context, input LoginInput) (string, error) {
    user, err := r.Service.Authenticate(input)
    if err != nil {
        return "", err
    }

    return GenerateToken(user.UserID, user.Role)
}

func (r *Resolver) FetchUsers(ctx context.Context) ([]*UserOutput, error) {
    users, err := r.Service.FetchUsers()
    if err != nil {
        return nil, err
    }

    var result []*UserOutput
    for _, u := range users {
        result = append(result, &UserOutput{
            UserID:   u.UserID,
            Username: u.Username,
            Email:    u.Email,
            Role:     u.Role,
        })
    }
    return result, nil
}
