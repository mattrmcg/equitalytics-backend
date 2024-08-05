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

	// Route to grab info about given ticker
	r.Get("/info/{ticker}", h.handleGetInfo)

}

// TODO: Define handler functions

func (h *InfoHandler) handlePing(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Pong!",
	})
}

func (h *InfoHandler) handleGetInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.service.GetInfoByTicker(r.Context(), chi.URLParam(r, "ticker")) // NEED TO FIX
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, info)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
