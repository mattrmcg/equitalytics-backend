package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

// TODO: Implement all three functions
// NEEDS TO BE DONE AFTER DATABASE IS SEEDED
func (userService *UserService) GetUserByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (userService *UserService) GetUserByID(id int) (*models.User, error) {
	return nil, nil
}

func (userService *UserService) CreateUser(user models.User) error {
	return nil
}
