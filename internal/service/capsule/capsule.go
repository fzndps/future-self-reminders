// Package service
package service

import (
	"context"

	"future-letter/internal/models"
)

type CapsuleService interface {
	CreateCapsule(ctx context.Context, userID int, input *models.CreateCapsuleInput) (*models.Capsule, error)
	GetCapsule(ctx context.Context, capsuleID, userID int) (*models.Capsule, error)
	GetUserCapsule(ctx context.Context, userID int) ([]models.Capsule, error)
	UpdateCapsule(ctx context.Context, capsuleID, userID int, input *models.UpdateCapsuleInput) (*models.Capsule, error)
	DeleteCapsule(ctx context.Context, capsuleID, userID int) error
	GetPendingForToday(ctx context.Context) ([]models.Capsule, error)
	MarkCapsulesAsSent(ctx context.Context, capsuleID int) error
}
