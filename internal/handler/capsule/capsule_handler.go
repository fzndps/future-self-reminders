// Package handler
package handler

import (
	"strconv"

	"future-letter/internal/middleware"
	"future-letter/internal/models"
	service "future-letter/internal/service/capsule"
	"future-letter/internal/utils"

	"github.com/gin-gonic/gin"
)

type CapsuleHandler struct {
	capsuleService service.CapsuleService
}

func NewCapsuleHandler(capsuleService service.CapsuleService) *CapsuleHandler {
	return &CapsuleHandler{
		capsuleService: capsuleService,
	}
}

// CreateCapsule handler
func (h *CapsuleHandler) CreateCapsule(c *gin.Context) {
	// mengambil user id dari middleware
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Bind input
	var input models.CreateCapsuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Panggil service dengan context
	capsule, err := h.capsuleService.CreateCapsule(c.Request.Context(), userID, &input)
	if err != nil {
		// Handle error yang berbeda
		errMsg := err.Error()
		if errMsg == "invalid date format, use YYYY-MM-DD" || errMsg == "due date must be in the future" {
			utils.BadRequestResponse(c, errMsg)
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to create capsule")
		return
	}

	utils.CreatedResponse(c, "Capsule created successfully", capsule.ToResponse())
}

func (h *CapsuleHandler) GetAllCapsules(c *gin.Context) {
	// dapatkan user ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.BadRequestResponse(c, "User not authenticated")
		return
	}

	// Panggil servce dengan context
	capsules, err := h.capsuleService.GetUserCapsule(c.Request.Context(), userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get capsules: "+err.Error())
		return
	}

	// Konversikan ke format respons
	responseCapsules := make([]*models.CapsuleResponse, 0, len(capsules))
	for i := range capsules {
		responseCapsules = append(responseCapsules, capsules[i].ToResponse())
	}

	utils.SuccessResponse(c, "Capsules retrieved successfully", responseCapsules)
}

func (h *CapsuleHandler) GetCapsuleByID(c *gin.Context) {
	// dapatkan userID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.BadRequestResponse(c, "User not authenticated")
		return
	}

	// dapatkan capsule ID
	capsuleIDStr := c.Param("capsuleID")
	capsuleID, err := strconv.Atoi(capsuleIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid capsule ID")
		return
	}

	// Panggil service dengan context
	capsule, err := h.capsuleService.GetCapsule(c.Request.Context(), capsuleID, userID)
	if err != nil {
		if err.Error() == "capsule not found" {
			utils.NotFoundResponse(c, "Capsule not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to get capsule: "+err.Error())
		return
	}

	utils.SuccessResponse(c, "Capsule retrieved successfully", capsule)
}

func (h *CapsuleHandler) UpdateCapsule(c *gin.Context) {
	// dapatkan user ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.BadRequestResponse(c, "User not authenticated")
		return
	}

	capsuleIDStr := c.Param("capsuleID")
	capsuleID, err := strconv.Atoi(capsuleIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid capsule ID")
		return
	}

	// Bind input
	var input models.UpdateCapsuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Panggil service dengan context
	capsule, err := h.capsuleService.UpdateCapsule(c.Request.Context(), capsuleID, userID, &input)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "cannot update capsule that is not pending" || errMsg == "invalid date format, use YYYY-MM-DD" || errMsg == "due date must be in the future" {
			utils.BadRequestResponse(c, errMsg)
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to update capsule: "+errMsg)
		return
	}

	// response
	utils.SuccessResponse(c, "Capsule updated successfully", capsule)
}

func (h *CapsuleHandler) DeleteCapsule(c *gin.Context) {
	// Dapatkan user ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.BadRequestResponse(c, "User not authenticated")
		return
	}

	// Dapatkan capsule ID
	capsuleIDStr := c.Param("capsuleID")
	capsuleID, err := strconv.Atoi(capsuleIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid capsule ID")
		return
	}

	// Panggil service dengan context
	err = h.capsuleService.DeleteCapsule(c.Request.Context(), capsuleID, userID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "capsule not found" {
			utils.NotFoundResponse(c, errMsg)
			return
		}

		if errMsg == "cannot delete capsule that is not pending" {
			utils.BadRequestResponse(c, errMsg)
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to delete capsule: "+errMsg)
		return
	}

	utils.SuccessResponse(c, "Capsule deleted successfully", nil)
}
