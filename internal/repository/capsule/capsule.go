// Package repository
package repository

import (
	"context"

	"future-letter/internal/models"
)

type CapsuleRepository interface {
	Create(ctx context.Context, capsule *models.Capsule) error
	GetByID(ctx context.Context, id int, userID int) (*models.Capsule, error)
	GetByUserID(ctx context.Context, userID int) ([]models.Capsule, error)
	Update(ctx context.Context, capsule *models.Capsule) error
	Delete(ctx context.Context, id, userID int) error
	GetPendingForToday(ctx context.Context) ([]models.Capsule, error)
	MarkAsSent(ctx context.Context, id int) error
}
