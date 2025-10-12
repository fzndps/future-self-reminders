package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"future-letter/internal/models"
	repository "future-letter/internal/repository/capsule"
)

type capsuleService struct {
	capsuleRepo repository.CapsuleRepository
}

func NewCapsuleService(capsuleRepo repository.CapsuleRepository) CapsuleService {
	return &capsuleService{
		capsuleRepo: capsuleRepo,
	}
}

var defaultStatus = "pending"

// CreateCapsule method untuk membuat capsule
func (s *capsuleService) CreateCapsule(ctx context.Context, userID int, input *models.CreateCapsuleInput) (*models.Capsule, error) {
	// Parse due date
	dueDate, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		return nil, errors.New("invalid format date, use YYYY-MM-DD")
	}

	now := time.Now()

	// Jika due date == hari ini, set jadi beberapa menit ke depan agar valid
	if dueDate.Format("2006-01-02") == now.Format("2006-01-02") {
		dueDate = now.Add(10 * time.Minute)
	}
	// Memastika due date di masa depan
	if dueDate.Before(now) {
		return nil, errors.New("due date must be in the future")
	}

	// buat object capsule
	capsule := &models.Capsule{
		UserID:         userID,
		Title:          input.Title,
		Message:        input.Message,
		DueDate:        dueDate,
		DeliveryMethod: input.DeliveryMethod,
		Status:         defaultStatus,
	}

	// handle optional fields
	if input.Category != "" {
		capsule.Category = sql.NullString{String: input.Category, Valid: true}
	}
	if input.Mood != "" {
		capsule.Mood = sql.NullString{String: input.Mood, Valid: true}
	}
	if input.ImageURL != "" {
		capsule.ImageURL = sql.NullString{String: input.ImageURL, Valid: true}
	}

	// Save to database
	err = s.capsuleRepo.Create(ctx, capsule)
	if err != nil {
		return nil, err
	}

	// Fetch full capsule
	fullCapsule, err := s.capsuleRepo.GetByID(ctx, capsule.ID, userID)
	if err != nil {
		return nil, err
	}

	return fullCapsule, nil
}

// GetCapsule mengambil kapsule berdasarkan id
func (s *capsuleService) GetCapsule(ctx context.Context, capsuleID, userID int) (*models.Capsule, error) {
	capsule, err := s.capsuleRepo.GetByID(ctx, capsuleID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get capsule: %v", err)
	}

	return capsule, nil
}

// GetUserCapsule mengambil capsule berdasarkan userID
func (s *capsuleService) GetUserCapsule(ctx context.Context, userID int) ([]models.Capsule, error) {
	capsule, err := s.capsuleRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get capsule by user id: %v", err)
	}

	return capsule, nil
}

func (s *capsuleService) UpdateCapsule(ctx context.Context, capsuleID, userID int, input *models.UpdateCapsuleInput) (*models.Capsule, error) {
	// ambil capsule yang ingin di update
	capsule, err := s.capsuleRepo.GetByID(ctx, capsuleID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get capsule: %v", err)
	}

	if capsule.Status != defaultStatus {
		return nil, errors.New("cannot update capsule that is not pending")
	}

	// Update fields jika kosong memakai yang lama
	if input.Title != "" {
		capsule.Title = input.Title
	}

	if input.Message != "" {
		capsule.Message = input.Message
	}

	if input.DeliveryMethod != "" {
		capsule.DeliveryMethod = input.DeliveryMethod
	}

	// Update due date jika di isi
	if input.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}

		if dueDate.Before(time.Now()) {
			return nil, errors.New("due date must be in the future")
		}
		capsule.DueDate = dueDate
	}

	if input.Category != "" {
		capsule.Category = sql.NullString{String: input.Category, Valid: true}
	}

	if input.Mood != "" {
		capsule.Mood = sql.NullString{String: input.Mood, Valid: true}
	}

	// Save update ke database
	err = s.capsuleRepo.Update(ctx, capsule)
	if err != nil {
		return nil, err
	}

	// Fetdh data terbaru dari database
	updateCapsule, err := s.capsuleRepo.GetByID(ctx, capsuleID, userID)
	if err != nil {
		return nil, err
	}

	return updateCapsule, nil
}

func (s *capsuleService) DeleteCapsule(ctx context.Context, capsuleID, userID int) error {
	// Get capsule yang mau dihapus
	capsule, err := s.capsuleRepo.GetByID(ctx, capsuleID, userID)
	if err != nil {
		return err
	}

	// cek status
	if capsule.Status != defaultStatus {
		return errors.New("cannot delete capsule that is not pending")
	}

	// Delete capsule
	return s.capsuleRepo.Delete(ctx, capsuleID, userID)
}

// Method ini akan digunakan oleh schedular
func (s *capsuleService) GetPendingCapsulesForToday(ctx context.Context) ([]models.Capsule, error) {
	return s.capsuleRepo.GetPendingForToday(ctx)
}

// Method ini dipanggil setelah email berhasil dikirim
func (s *capsuleService) MarkCapsulesAsSent(ctx context.Context, capsuleID int) error {
	return s.capsuleRepo.MarkAsSent(ctx, capsuleID)
}
