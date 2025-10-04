// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "future-letter/internal/models/user"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Create untuk input user ke database
func (u *userRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (name, email, password, timezone) VALUES (?, ?, ?, ?)"

	result, err := u.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.Timezone,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id user: %w", err)
	}

	user.ID = int(id)

	return nil
}

// GetByID untuk mendapatkan user berdasarkan ID
func (u *userRepositoryImpl) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}

	query := "SELECT id, name, email, password, timezone, created_at, updated_at FROM users WHERE id = ?"

	err := u.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Timezone, &user.CreatedAt, &user.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail untuk mendapatkan user berdasarkan email
func (u *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	query := "SELECT id, name, email, password, timezone, created_at, updated_at FROM users WHERE email = ?"

	err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Timezone, &user.CreatedAt, &user.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update untuk mengubah nama dan timezone
func (u *userRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET name = ?, timezone = ? WHERE id = ?"

	_, err := u.db.ExecContext(ctx, query, user.Name, user.Timezone, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete untuk menghapus user yang ada di database
func (u *userRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
