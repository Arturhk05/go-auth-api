package services

import (
	"fmt"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/internal/models"
	"github.com/arturhk05/go-auth-api/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *AuthService) Register(password string, email string, username string) (uuid.UUID, error) {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return uuid.Nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.Security.BcryptCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

func (s *AuthService) Login(email string, password string) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// TODO: Implement account lockout after multiple failed attempts
	// TODO: Email verification check

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// TODO: Generate JWT token

	return &models.LoginResponse{
		UserID:      user.ID,
		AccessToken: "token",
	}, nil
}
