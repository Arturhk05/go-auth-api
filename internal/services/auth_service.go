package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arturhk05/go-auth-api/config"
	apperrors "github.com/arturhk05/go-auth-api/internal/errors"
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
		return nil, fmt.Errorf("register: %w", apperrors.ErrUserAlreadyExists)
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, fmt.Errorf("register: check user existence: %w", err)
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
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	// TODO: Email verification check

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
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, fmt.Errorf("login: %w", apperrors.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("login: get user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("login: %w", apperrors.ErrAccountInactive)
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now().UTC()) {
		return nil, fmt.Errorf("login: %w", apperrors.ErrAccountLocked)
	}

	if user.FailedLoginAttempts >= s.cfg.Security.MaxLoginAttempts {
		lockErr := s.userRepo.LockAccountAndResetLoginAttempts(user.ID, sql.NullTime{Valid: true, Time: time.Now().UTC().Add(time.Duration(s.cfg.Security.LockDurationMinutes) * time.Minute)})
		if lockErr != nil {
			return nil, fmt.Errorf("login: lock account: %w", lockErr)
		}
		return nil, fmt.Errorf("login: %w", apperrors.ErrAccountLocked)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		updateErr := s.userRepo.UpdateLoginAttempts(user.ID, user.FailedLoginAttempts+1)
		if updateErr != nil {
			return nil, fmt.Errorf("login: update login attempts: %w", updateErr)
		}
		return nil, fmt.Errorf("login: %w", apperrors.ErrInvalidCredentials)
	}

	// All creadentials are valid from this point

	err = s.refreshTokenRepo.RevokeByUserId(user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: revoke tokens: %w", err)
	}
	err = s.userRepo.ResetLoginAttempts(user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: reset login attempts: %w", err)
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.ExpirationHours)
	if err != nil {
		return nil, fmt.Errorf("login: generate access token: %w", err)
	}

	refreshToken, err := s.generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: generate refresh token: %w", err)
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
		return nil, fmt.Errorf("refresh token: validate: %w", err)
	}

	user, err := s.userRepo.GetUserById(refreshClaims.UserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, fmt.Errorf("refresh token: %w", apperrors.ErrUserNotFound)
		}
		return nil, fmt.Errorf("refresh token: get user: %w", err)
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now().UTC()) {
		return nil, fmt.Errorf("refresh token: %w", apperrors.ErrAccountLocked)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("refresh token: %w", apperrors.ErrAccountInactive)
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

func (s *AuthService) validateRefreshToken(refreshToken string) (*utils.RefreshClaims, error) {
	refreshClaims, err := utils.ValidateRefreshToken(refreshToken, s.cfg.JWT.RefreshSecret)
	if err != nil {
		if errors.Is(err, apperrors.ErrTokenExpired) {
			return nil, fmt.Errorf("validate refresh token: %w", apperrors.ErrTokenExpired)
		}
		if errors.Is(err, apperrors.ErrInvalidToken) {
			return nil, fmt.Errorf("validate refresh token: %w", apperrors.ErrInvalidToken)
		}
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	refreshTokenHash := utils.HashToken(refreshToken)
	_, err = s.refreshTokenRepo.ValidateRefreshToken(refreshTokenHash)
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", apperrors.ErrTokenRevoked)
	}

	return refreshClaims, nil
}

func (s *AuthService) generateAndSaveRefreshToken(userID uuid.UUID) (string, error) {
	refreshToken, err := utils.GenerateRefreshToken(userID, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshExpirationHours)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: generate: %w", err)
	}
	refreshTokenHashed := utils.HashToken(refreshToken)
	refreshTokenExpirationTime := time.Now().UTC().Add(time.Duration(s.cfg.JWT.RefreshExpirationHours) * time.Hour)

	err = s.refreshTokenRepo.Create(userID, refreshTokenHashed, refreshTokenExpirationTime)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: save: %w", err)
	}

	return refreshToken, nil
}
