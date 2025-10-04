// Package repository
package repository

import (
	"context"
	"database/sql"

	models "future-letter/internal/models/user"
)

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (u *UserRepositoryImpl) Created(ctx context.Context, user *models.User) error {
}

func (u *UserRepositoryImpl) FindUserByID(ctx context.Context, id int) (*models.User, error) {
}

func (u *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
}
