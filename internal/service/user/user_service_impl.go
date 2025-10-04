// Package service
package service

import (
	"context"

	models "future-letter/internal/models/user"
	repository "future-letter/internal/repository/user"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, input *models.RegisterInput) (*models.User, error) {
}

func (s *userService) Login(ctx context.Context, input *models.LoginInput) (*models.User, error) {
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, input *models.UpdateProfileInput) (*models.User, error) {
}

func (s *userService) DeleteAccount(ctx context.Context, userID int) error {
}
