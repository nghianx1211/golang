package user

import (
    "context"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v5"
    "time"
	"golang/internal/model"
)

var jwtKey = []byte("secret") // Đổi thành env trong thực tế

type Resolver struct {
    DB *gorm.DB
}

type Claims struct {
    UserID uint
    Role   string
    jwt.RegisteredClaims
}

func (r *Resolver) CreateUser(ctx context.Context, username, email, password, role string) (*User, error) {
    if role != "manager" && role != "member" {
        return nil, errors.New("invalid role")
    }

    hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    user := model.User{
        Username:     username,
        Email:        email,
        Role:         role,
        PasswordHash: string(hashed),
    }

    if err := r.DB.Create(&user).Error; err != nil {
        return nil, errors.New("email already exists")
    }

    return &user, nil
}

func (r *Resolver) Login(ctx context.Context, email, password string) (string, *model.User, error) {
    var user model.User
    if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    claims := &Claims{
        UserID: user.UserID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(jwtKey)

    return tokenString, &user, nil
}

func (r *Resolver) FetchUsers(ctx context.Context) ([]model.User, error) {
    claims := ctx.Value("claims").(*Claims)
    if claims == nil {
        return nil, errors.New("unauthorized")
    }

    var users []model.User
    r.DB.Find(&users)
    return users, nil
}

func (r *Resolver) Logout(ctx context.Context) (bool, error) {
    // Nếu dùng JWT thì logout chỉ là xóa client token
    return true, nil
}
