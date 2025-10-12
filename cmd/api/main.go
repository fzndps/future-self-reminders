package main

import (
	"log"

	"future-letter/internal/config"
	"future-letter/internal/database"
	capsuleRepository "future-letter/internal/repository/capsule"
	userRepository "future-letter/internal/repository/user"
	"future-letter/internal/routes"
	capsuleService "future-letter/internal/service/capsule"
	emailService "future-letter/internal/service/email"
	schedulerService "future-letter/internal/service/scheduler"
	userService "future-letter/internal/service/user"
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
	userRepo := userRepository.NewUserRepository(database.DB)
	capsuleRepo := capsuleRepository.NewCapsuleRepository(database.DB)

	// Initalize service
	userSvc := userService.NewUserService(userRepo)
	capsuleSvc := capsuleService.NewCapsuleService(capsuleRepo)
	emailSvc := emailService.NewEmailService(cfg)

	// Scheduler service
	schedulerSvc := schedulerService.NewSchedulerService(cfg, userRepo, capsuleSvc, emailSvc)
	err = schedulerSvc.Start()
	if err != nil {
		log.Fatal("failed to start scheduler:", err)
	}

	defer schedulerSvc.Stop()

	// Setup routes
	routes.SetupRoutes(router, cfg, userSvc, capsuleSvc)

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
