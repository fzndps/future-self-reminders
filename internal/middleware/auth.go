// Package middleware berfungsi menjalankan middleware sebelum handler
package middleware

import (
	"strings"

	"future-letter/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	// Return function akan dijalankan saat ada request
	return func(c *gin.Context) {
		// authHeader ini bisasanya berbentuk : Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Split berfungsi untuk mengubah "Bearer <token>" menjadi ["Bearer", "<token>"]
		parts := strings.Split(authHeader, " ")

		// Validasi format harus ada 2 parts dan parts pertama harus berisi "Bearer"
		if len(parts) != 2 || parts[0] == "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization format.")
			c.Abort()
			return
		}

		// Parts ke 2 adalah token
		tokenString := parts[1]

		// ValidateJWT dari utils akan mengecek token masih valid dan belum expired
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// c.Set untuk menyimpan data ke context dan bisa di ambil
		// dihandler untuk mengetahui siapa yang login dengan c.Get
		c.Set("userID", claims.ID)
		c.Set("email", claims.Email)

		// Lanjut ke middleware/handler berikutnya
		c.Next()
	}
}

// GetUserID untuk mengambil userID
func GetUserID(c *gin.Context) (int, bool) {
	// Ambil data dari context
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	// Convert nilai interface/any ke int
	id, ok := userID.(int)
	if !ok {
		return 0, false
	}

	return id, true
}

func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}

	emailSTR, ok := email.(string)
	if !ok {
		return "", false
	}

	return emailSTR, true
}
