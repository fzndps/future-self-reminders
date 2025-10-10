// Package handler
package handler

import (
	"future-letter/internal/config"
	service "future-letter/internal/service/capsule"
)

type capsuleHandler struct {
	capsuleService service.CapsuleService
	cfg            *config.Config
}

func NewCapsuleHandler(capsuleService service.CapsuleService, cfg *config.Config) *capsuleHandler {
	return &capsuleHandler{
		capsuleService: capsuleService,
		cfg:            cfg,
	}
}
