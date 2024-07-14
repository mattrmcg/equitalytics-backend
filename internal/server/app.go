package server

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mattrmcg/equitalytics-backend/api"
	"github.com/mattrmcg/equitalytics-backend/api/handlers"
	"github.com/mattrmcg/equitalytics-backend/internal/services/info"
	"github.com/mattrmcg/equitalytics-backend/internal/services/user"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

func NewAPIServer(addr string, dbPool *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr: addr,
		db:   dbPool,
	}
}

func (s *APIServer) Run() error {
	userService := user.NewUserService(s.db)
	userHandler := handlers.NewUserHandler(userService)

	infoService := info.NewInfoService(s.db)
	infoHandler := handlers.NewInfoHandler(infoService)

	router := api.SetupRouter(userHandler, infoHandler)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
