package user

type CreateUserInput struct {
    Username string
    Email    string
    Password string
    Role     string
}

type LoginInput struct {
    Email    string
    Password string
}

type UserOutput struct {
    UserID   string
    Username string
    Email    string
    Role     string
}
