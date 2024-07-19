package info

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

type InfoService struct {
	db *pgxpool.Pool
}

func NewInfoService(db *pgxpool.Pool) *InfoService {
	return &InfoService{db: db}
}

func (infoService *InfoService) GetInfoByCIK(cik int) (*models.CompanyInfo, error) {
	return nil, nil
}
