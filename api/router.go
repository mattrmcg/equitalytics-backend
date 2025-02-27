package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/mattrmcg/equitalytics-backend/api/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler, infoHandler *handlers.InfoHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(httprate.LimitAll(50, time.Second))
	r.Use(httprate.LimitByIP(30, time.Minute))
	r.Use(middleware.Timeout(10 * time.Second))

	userHandler.RegisterRoutes(r)
	infoHandler.RegisterRoutes(r)

	return r
}
