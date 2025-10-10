// Package handler
package handler

import (
	"future-letter/internal/config"
	"future-letter/internal/middleware"
	"future-letter/internal/models"
	service "future-letter/internal/service/user"
	"future-letter/internal/utils"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	userService service.UserService
	cfg         *config.Config
}

func NewUserHandler(userService service.UserService, cfg *config.Config) *authHandler {
	return &authHandler{
		userService: userService,
		cfg:         cfg,
	}
}

// Handler Register untuk menangani request dan response Register
func (h *authHandler) Register(c *gin.Context) {
	var input models.RegisterInput

	// Bind request body
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Panggil service dengan context
	user, err := h.userService.Register(c.Request.Context(), &input)
	if err != nil {
		if err.Error() == "email already registered" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to register")
		return
	}

	// Generate jwt token
	token, err := utils.GenerateToken(user.ID, user.Email, h.cfg.JWT.Expiry)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token")
		return
	}

	responseData := map[string]any{
		"user":  user.ToResponse(),
		"token": token,
	}

	utils.CreatedResponse(c, "User registered successfully", responseData)
}

func (h *authHandler) Login(c *gin.Context) {
	var input models.LoginInput

	// Bind request body
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Panggil service dengan context
	user, err := h.userService.Login(c.Request.Context(), &input)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, h.cfg.JWT.Expiry)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	responseData := map[string]any{
		"user":  user.ToResponse(),
		"token": token,
	}

	utils.SuccessResponse(c, "Login successfully", responseData)
}

func (h *authHandler) GetProfile(c *gin.Context) {
	// Get user id dari context (middleware)
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Panggil service dengan context
	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", user.ToResponse())
}

func (h *authHandler) UpdateProfile(c *gin.Context) {
	// Get user id dari context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var input models.UpdateProfileInput

	// Bind input body
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Panggil service dengan context
	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &input)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile")
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", user.ToResponse())
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	// Get user id dari context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Get user email dari context
	email, ok := middleware.GetEmail(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Generate token baru
	newToken, err := utils.GenerateToken(userID, email, h.cfg.JWT.Expiry)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate new token")
		return
	}

	resposeData := map[string]any{
		"token": newToken,
	}

	utils.SuccessResponse(c, "Token refreshed successfully", resposeData)
}
