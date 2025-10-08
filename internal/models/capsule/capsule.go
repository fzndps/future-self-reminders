// Package models
package models

import (
	"database/sql"
	"time"
)

// Capsule struct
type Capsule struct {
	ID             int            `json:"id" db:"id"`
	UserID         int            `json:"user_id" db:"user_id"`
	Title          string         `json:"title" db:"title"`
	Message        string         `json:"message" db:"message"`
	DueDate        time.Time      `json:"due_date" db:"due_date"`
	DeliveryMethod string         `json:"delivery_method" db:"delivery_method"`
	Status         string         `json:"status" db:"status"`
	Category       sql.NullString `json:"category" db:"category"`
	Mood           sql.NullString `json:"mood" db:"mood"`
	ImageURL       sql.NullString `json:"image_url" db:"image_url"`
	SentAt         sql.NullTime   `json:"sent_at" db:"sent_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

// DTO

// CreateCapsuleInput DTO untuk membuat capsule baru
type CreateCapsuleInput struct {
	Title          string `json:"title" binding:"required"`
	Message        string `json:"message" binding:"required"`
	DueDate        string `json:"due_date" binding:"required"`
	DeliveryMethod string `json:"delivery_method" binding:"required"`
	Status         string `json:"status"`
	Category       string `json:"category"`
	Mood           string `json:"mood" `
	ImageURL       string `json:"image_url"`
}

// UpdateCapsuleInput DTO untuk mengupdate capsule
type UpdateCapsuleInput struct {
	Title          string `json:"title"`
	Message        string `json:"message"`
	DueDate        string `json:"due_date"`
	DeliveryMethod string `json:"delivery_method"`
	Status         string `json:"status"`
	Category       string `json:"category"`
	Mood           string `json:"mood" `
}

type CapsuleResponse struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	Title          string     `json:"title"`
	Message        string     `json:"message"`
	DueDate        string     `json:"due_date"`
	DeliveryMethod string     `json:"delivery_method"`
	Status         string     `json:"status"`
	Category       *string    `json:"category"`
	Mood           *string    `json:"mood"`
	ImageURL       *string    `json:"image_url"`
	SentAt         *time.Time `json:"sent_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ToResponse mengkonversi capsule ke CapsuleResponse
func (c *Capsule) ToResponse() *CapsuleResponse {
	response := &CapsuleResponse{
		ID:             c.ID,
		UserID:         c.UserID,
		Title:          c.Title,
		Message:        c.Message,
		DeliveryMethod: c.DeliveryMethod,
		Status:         c.Status,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}

	// Handle nullable fields
	if c.Category.Valid {
		response.Category = &c.Category.String
	}
	if c.Mood.Valid {
		response.Mood = &c.Mood.String
	}
	if c.ImageURL.Valid {
		response.ImageURL = &c.ImageURL.String
	}
	if c.SentAt.Valid {
		response.SentAt = &c.SentAt.Time
	}

	return response
}
