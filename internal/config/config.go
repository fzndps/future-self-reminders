// Package config untuk menampung semua konfigurasi aplikasi yang dibuat
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database  DatabaseConfig
	App       AppConfig
	JWT       JWTConfig
	Email     EmailConfig
	Schedular SchedularConfig
}

// DatabaseConfig menampung konfigurasi database MYSQL
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// AppConfig menampung konfigurasi aplikasi
type AppConfig struct {
	Port string
	Env  string
}

// JWTConfig menampung konfigurasi JWT
type JWTConfig struct {
	Secret string
	Expiry int
}

// EmailConfig menampung konfigurasi email SMTP
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

// SchedularConfig menampung konfigurasi schedular
type SchedularConfig struct {
	CronExpression string
	Timezone       string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},

		App: AppConfig{
			Port: os.Getenv("APP_PORT"),
			Env:  os.Getenv("APP_ENV"),
		},

		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Expiry: getENVasInt("JWT_EXPIRY", 72),
		},

		Email: EmailConfig{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     getENVasInt("SMTP_PORT", 587),
			SMTPUsername: os.Getenv("SMTP_USERNAME"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			SMTPFrom:     os.Getenv("SMTP_FROM"),
		},

		Schedular: SchedularConfig{
			CronExpression: os.Getenv("SCHEDULER_CRON"),
			Timezone:       os.Getenv("SCHEDULER_TIMEZONE"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	// Cek database config
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Cek JWT secret (sangat penting untuk keamanan!)
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	// Cek app port
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	// Cek email config (jika ingin kirim email)
	if c.Email.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.Email.SMTPUsername == "" {
		return fmt.Errorf("SMTP_USERNAME is required")
	}
	if c.Email.SMTPPassword == "" {
		return fmt.Errorf("SMTP_PASSWORD is required")
	}

	// Semua validasi passed
	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

func getENVasInt(key string, defaultValue int) int {
	valueSTR := os.Getenv(key)

	if valueSTR == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueSTR)
	if err != nil {
		return defaultValue
	}

	return value
}
