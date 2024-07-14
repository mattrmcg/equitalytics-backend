package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

type InfoHandler struct {
	service models.InfoService
}

func NewInfoHandler(service models.InfoService) *InfoHandler {
	return &InfoHandler{service: service}
}

func (h *InfoHandler) RegisterRoutes(router chi.Router) {
	// TODO: Add routes to router

}

// TODO: Define handler functions
