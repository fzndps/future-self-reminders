// Package routes untuk mengkonfigurasi semua routing endpoint
package routes

import (
	"future-letter/internal/config"
	"future-letter/internal/database"
	handler "future-letter/internal/handler/user"
	service "future-letter/internal/service/user"
	"future-letter/internal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, userService service.UserService) {
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
		authHandler := handler.NewUserHandler(userService, cfg)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Protected endpoints
			auth.GET("/profile", authHandler.GetProfile)
			auth.PUT("/update", authHandler.UpdateProfile)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
	}
}
