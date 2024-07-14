package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mattrmcg/equitalytics-backend/api/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler, infoHandler *handlers.InfoHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userHandler.RegisterRoutes(r)
	infoHandler.RegisterRoutes(r)

	return r
}
