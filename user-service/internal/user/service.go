package user

import (
	"errors"

	"github.com/nghianx1211/golang/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateUser(input CreateUserInput) (*model.User, error) {
	if input.Role != "manager" && input.Role != "member" {
		return nil, errors.New("invalid role")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username:     input.Username,
		Email:        input.Email,
		Role:         input.Role,
		PasswordHash: string(hash),
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Authenticate(input LoginInput) (*model.User, error) {
	var user model.User
	if err := s.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *Service) FetchUsers() ([]model.User, error) {
	var users []model.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
