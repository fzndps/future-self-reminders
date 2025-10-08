package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "future-letter/internal/models/capsule"
)

type capsuleRepository struct {
	db *sql.DB
}

func NewCapsuleRepository(db *sql.DB) CapsuleRepository {
	return &capsuleRepository{
		db: db,
	}
}

// Create
func (r *capsuleRepository) Create(ctx context.Context, capsule *models.Capsule) error {
	query := "INSERT INTO capsules (user_id, title, message, due_date, delivery_method, category, mood, image_url, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, capsule.UserID, capsule.Title, capsule.Message, capsule.DueDate, capsule.DeliveryMethod, capsule.Category, capsule.Mood, capsule.ImageURL, capsule.Status)
	if err != nil {
		return fmt.Errorf("failed to create capsule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	capsule.ID = int(id)
	return nil
}

func (r capsuleRepository) GetByID(ctx context.Context, id int, userID int) (*models.Capsule, error) {
	capsule := &models.Capsule{}
	query := `SELECT 
		id, user_id, title, message, due_date, delivery_method, status, category, mood, image_url, sent_at, created_at, updated_at
		FROM capsules
		WHERE id = ? AND user_id = ?
	`

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&capsule.ID,
		&capsule.UserID,
		&capsule.Title,
		&capsule.Message,
		&capsule.DueDate,
		&capsule.DeliveryMethod,
		&capsule.Status,
		&capsule.Category,
		&capsule.Mood,
		&capsule.ImageURL,
		&capsule.SentAt,
		&capsule.CreatedAt,
		&capsule.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("capsule not found")
		}
		return nil, err
	}

	return capsule, nil
}

func (r *capsuleRepository) GetByUserID(ctx context.Context, userID int) ([]models.Capsule, error) {
	query := `SELECT 
		id, user_id, title, message, due_date, delivery_method, status, category, mood, image_url, sent_at, created_at, updated_at
		FROM capsules
		WHERE user_id = ?
		ORDER BY due_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get capsules by id: %w", err)
	}

	defer rows.Close()

	capsules := []models.Capsule{}

	for rows.Next() {
		var capsule models.Capsule
		err := rows.Scan(
			&capsule.ID,
			&capsule.UserID,
			&capsule.Title,
			&capsule.Message,
			&capsule.DueDate,
			&capsule.DeliveryMethod,
			&capsule.Status,
			&capsule.Category,
			&capsule.Mood,
			&capsule.ImageURL,
			&capsule.SentAt,
			&capsule.CreatedAt,
			&capsule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		capsules = append(capsules, capsule)
	}
	return capsules, nil
}

func (r *capsuleRepository) Update(ctx context.Context, capsule *models.Capsule) error {
	query := `UPDATE capsules
		SET title = ?, message = ?, due_date = ?, delivery_method = ?, category = ?, mood = ?
		WHERE id = ? AND user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, capsule.Title, capsule.Message, capsule.DueDate, capsule.DeliveryMethod, nullIfEmpty(capsule.Category.String), nullIfEmpty(capsule.Mood.String), capsule.ID, capsule.UserID)

	return err
}

func (r *capsuleRepository) Delete(ctx context.Context, id, userID int) error {
	query := "UPATE capsules SET status = 'canceled' WHERE id = ? AND user_id = ?"

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return errors.New("capsules not found")
	}

	return nil
}

func (r *capsuleRepository) GetPendingForToday(ctx context.Context) ([]models.Capsule, error) {
	query := `
		SELECT id, user_id, title, message, due_date, delivery_method,
			status, category, mood, image_url, sent_at, created_at, updated_at
		FROM capsules
		WHERE DATE(due_date) = CUREDATE() AND status = 'pending'
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	defer rows.Close()

	capsules := []models.Capsule{}
	for rows.Next() {
		var capsule models.Capsule
		err := rows.Scan(
			&capsule.ID,
			&capsule.UserID,
			&capsule.Title,
			&capsule.Message,
			&capsule.DueDate,
			&capsule.DeliveryMethod,
			&capsule.Status,
			&capsule.Category,
			&capsule.Mood,
			&capsule.ImageURL,
			&capsule.SentAt,
			&capsule.CreatedAt,
			&capsule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		capsules = append(capsules, capsule)
	}

	return capsules, nil
}

func (r *capsuleRepository) MarkAsSent(ctx context.Context, id int) error {
	query := `UPDATE capsules 
		SET status = 'sent', sent_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
