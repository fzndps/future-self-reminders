// Package database
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"future-letter/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("error connecting to database: %w", err)
	}

	DB = db

	log.Println("Successfully connected to mysql database")

	return nil
}

func CloseDB() error {
	if DB != nil {
		if err := DB.Close(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}

		log.Println("Database connection closed")
	}

	return nil
}

func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
