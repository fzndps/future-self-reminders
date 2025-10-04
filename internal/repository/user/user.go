// Package repository
package repository

import (
	"context"

	models "future-letter/internal/models/user"
)

type UserRepository interface {
	Created(ctx context.Context, user *models.User) error
	FindUserByID(ctx context.Context, id int) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}
