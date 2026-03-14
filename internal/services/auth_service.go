package services

import (
	"fmt"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/internal/models"
	"github.com/arturhk05/go-auth-api/internal/repositories"
	"github.com/arturhk05/go-auth-api/internal/utils"
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

func (s *AuthService) Register(password string, email string, username string) (*models.AuthResponse, error) {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.Security.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.AuthResponse{
		UserID:      user.ID,
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) Login(email string, password string) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// TODO: Implement account lockout after multiple failed attempts
	// TODO: Email verification check
	// TODO: Verify locked account status

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.AuthResponse{
		UserID:      user.ID,
		AccessToken: accessToken,
	}, nil
}
