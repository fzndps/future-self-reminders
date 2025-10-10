// Package service
package service

import (
	"context"

	"future-letter/internal/models"
)

type UserService interface {
	Register(ctx context.Context, input *models.RegisterInput) (*models.User, error)
	Login(ctx context.Context, input *models.LoginInput) (*models.User, error)
	GetProfile(ctx context.Context, userID int) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int, input *models.UpdateProfileInput) (*models.User, error)
	DeleteAccount(ctx context.Context, userID int) error
}
