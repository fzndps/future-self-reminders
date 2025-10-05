package main

import (
	"log"

	"future-letter/internal/config"
	"future-letter/internal/database"
	repository "future-letter/internal/repository/user"
	"future-letter/internal/routes"
	service "future-letter/internal/service/user"
	"future-letter/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Initalize database
	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	defer database.CloseDB()

	// Initalize jwt
	utils.InitJWT(cfg.JWT.Secret)

	// Setup gin router
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Initalize repository
	userRepo := repository.NewUserRepository(database.DB)

	// Initalize service
	userService := service.NewUserService(userRepo)

	// Setup routes
	routes.SetupRoutes(router, cfg, userService)

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
