package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

type UserHandler struct {
	service models.UserService
}

func NewUserHandler(service models.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(router chi.Router) {
	// TODO: Add routes
}

// TODO: Define handler functions
