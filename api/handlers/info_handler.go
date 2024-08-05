package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
	"github.com/mattrmcg/equitalytics-backend/pkg/utils"
)

type InfoHandler struct {
	service models.InfoService
}

func NewInfoHandler(service models.InfoService) *InfoHandler {
	return &InfoHandler{service: service}
}

func (h *InfoHandler) RegisterRoutes(r chi.Router) {
	// TODO: Add routes to router
	r.Get("/ping", h.handlePing) // Auxiliary ping route
	// r.Get("/companyinfo/{ticker}", getCompanyInfo())

}

// TODO: Define handler functions

func (h *InfoHandler) handlePing(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Pong!",
	})
}
