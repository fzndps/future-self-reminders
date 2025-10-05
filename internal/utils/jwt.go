package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtsecret []byte

func InitJWT(secret string) {
	jwtsecret = []byte(secret)
}

// GenerateToken membuat JWT token baru untuk user
func GenerateToken(userID int, email string, expiryHours int) (string, error) {
	// Cek jika jwt sudah diinisialisasi
	if len(jwtsecret) == 0 {
		return "", errors.New("JWT secret not initialize")
	}

	// Buat claims / data yang disimpan di token
	claims := &JWTClaims{
		ID:    userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt -> kapan token expire
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiryHours))),
			// IssuedAt -> kapan token dibuat
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// NotBefore -> waktu token mulai valid
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// buat token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Encode token dan tanda tangan yang hasilnya jwt yang bisa dikirim ke client
	tokenString, err := token.SignedString(jwtsecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT memvalidasi JWT token dan extract claims
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	// Cek apakah jwt secret sudah di inisialisasi
	if len(jwtsecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	// Parse token string dam validasi signature
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return jwtsecret, nil
	})
	// Cek bila ada error saat parse
	if err != nil {
		return nil, err
	}

	// Extract claims -> ambil claims dari token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Validasi token cek apakah token valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func RefreshToken(oldTokenString string, expiryHours int) (string, error) {
	// validasi token lama
	claims, err := ValidateJWT(oldTokenString)
	if err != nil {
		return "", err
	}

	// Generate token baru dengan claims yang sama
	// gunakan expiry time baru
	return GenerateToken(claims.ID, claims.Email, expiryHours)
}

func ExtractUserIDFromToken(tokenString string) (int, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.ID, nil
}
