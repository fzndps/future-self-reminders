// Package routes untuk mengkonfigurasi semua routing endpoint
package routes

import (
	"future-letter/internal/config"
	"future-letter/internal/database"
	capsuleHandler "future-letter/internal/handler/capsule"
	userHandler "future-letter/internal/handler/user"
	"future-letter/internal/middleware"
	capsuleService "future-letter/internal/service/capsule"
	userService "future-letter/internal/service/user"
	"future-letter/internal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, userService userService.UserService, capsuleService capsuleService.CapsuleService) {
	// health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		// cek database health
		err := database.HealthCheck()
		if err != nil {
			utils.InternalServerErrorResponse(c, "Database connection failed")
			return
		}

		utils.SuccessResponse(c, "API is healtht", map[string]any{
			"status":   "ok",
			"database": "connected",
		})
	})

	api := router.Group("/api/v1")
	{
		// Auth routes
		// Initialize auth handler dengan dependency injection
		authHandler := userHandler.NewUserHandler(userService, cfg)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Protected endpoints
			auth.GET("/profile", middleware.AuthRequired(), authHandler.GetProfile)
			auth.PUT("/update", middleware.AuthRequired(), authHandler.UpdateProfile)
			auth.POST("/refresh", middleware.AuthRequired(), authHandler.RefreshToken)
		}

		// Initialize capsule hadnler dengan dependency injection
		capsuleHandler := capsuleHandler.NewCapsuleHandler(capsuleService)

		capsules := api.Group("/capsules")
		capsules.Use(middleware.AuthRequired())
		{
			capsules.GET("", capsuleHandler.GetAllCapsules)
			capsules.POST("", capsuleHandler.CreateCapsule)
			capsules.GET("/:capsuleID", capsuleHandler.GetCapsuleByID)
			capsules.PUT("/:capsuleID", capsuleHandler.UpdateCapsule)
			capsules.DELETE("/:capsuleID", capsuleHandler.DeleteCapsule)
		}
	}
}
