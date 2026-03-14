package services

import (
	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/internal/models"
	"github.com/arturhk05/go-auth-api/internal/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserById(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()

	return response, nil
}
