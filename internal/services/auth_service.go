package services

import (
	"fmt"
	"time"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/internal/models"
	"github.com/arturhk05/go-auth-api/internal/repositories"
	"github.com/arturhk05/go-auth-api/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	cfg              *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, refreshTokenRepo *repositories.RefreshTokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
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

	refreshToken, err := s.generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	s.refreshTokenRepo.RevokeByUserId(user.ID)

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	refreshClaims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserById(refreshClaims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// TODO: Verify locked account status

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.AuthResponse{
		UserID:      user.ID,
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) validateRefreshToken(refreshToken string) (*utils.RefreshClaims, error) {
	refreshClaims, err := utils.ValidateRefreshToken(refreshToken, s.cfg.JWT.RefreshSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := utils.HashToken(refreshToken)
	_, err = s.refreshTokenRepo.ValidateRefreshToken(refreshTokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token is revoked or invalid: %w", err)
	}

	return refreshClaims, nil
}

func (s *AuthService) generateAndSaveRefreshToken(userID uuid.UUID) (string, error) {
	refreshToken, err := utils.GenerateRefreshToken(userID, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshExpirationHours)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenHashed := utils.HashToken(refreshToken)
	refreshTokenExpirationTime := time.Now().UTC().Add(time.Duration(s.cfg.JWT.RefreshExpirationHours) * time.Hour)

	err = s.refreshTokenRepo.Create(userID, refreshTokenHashed, refreshTokenExpirationTime)
	if err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return refreshToken, nil
}
