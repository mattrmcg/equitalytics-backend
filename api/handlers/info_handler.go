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

	// Route to grab all tickers in database
	r.Get("/info/tickers", h.handleGetTickers)

}

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

func (h *InfoHandler) handleGetTickers(w http.ResponseWriter, r *http.Request) {
	tickers, err := h.service.GetAllTickers(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	type tickerValue struct {
		Value string `json:"value"`
	}

	var tickerList []tickerValue
	for _, ticker := range tickers {
		tickerList = append(tickerList, tickerValue{Value: ticker})
	}

	err = utils.WriteJSON(w, http.StatusOK, tickerList)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
