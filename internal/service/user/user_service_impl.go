// Package service
package service

import (
	"context"
	"errors"
	"fmt"

	models "future-letter/internal/models/user"
	repository "future-letter/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register untuk mendaftarkan user baru
func (s *userService) Register(ctx context.Context, input *models.RegisterInput) (*models.User, error) {
	// Cek email apakah sudah ada
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password agar aman
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Jika timezone kosong maka akan di isi dengan Asia/Jakarta
	timezone := input.Timezone
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	// Buat objek user
	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPass),
		Timezone: timezone,
	}

	// Simpan user ke database
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Dapatkan user yang sudah masuk
	fullUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}

	return fullUser, nil
}

// Login untuk masuk ke app
func (s *userService) Login(ctx context.Context, input *models.LoginInput) (*models.User, error) {
	// Dapatkan user berdasarkan email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetProfile untuk mendapat informasi profil
func (s *userService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %v", userID, err)
	}

	return user, nil
}

// UpdateProfile untuk mengupdate nama dan timezone
func (s *userService) UpdateProfile(ctx context.Context, userID int, input *models.UpdateProfileInput) (*models.User, error) {
	// Cari user yang ingin di update berdasarkan id
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %v", userID, err)
	}

	// Ubah data
	user.Name = input.Name
	user.Timezone = input.Timezone

	// Save ke database
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed save to database: %v", err)
	}

	// Fetch data yang sudah di update
	userUpdate, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userUpdate, nil
}

// DeleteAccount menghapus akun yang sudah ada
func (s *userService) DeleteAccount(ctx context.Context, userID int) error {
	return s.userRepo.Delete(ctx, userID)
}
