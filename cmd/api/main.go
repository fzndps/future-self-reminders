package main

import (
	"log"

	"future-letter/internal/config"
	"future-letter/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	defer database.CloseDB()
}
