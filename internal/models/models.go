package models

import "time"

type UserService interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type InfoService interface {
	GetInfoByCIK(cik int) (*Info, error)
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	APIKey    string    `json:"-"`       // figure out optimal way of storing api keys
	Allowed   bool      `json:"allowed"` // allow access to API
	CreatedAt time.Time `json:"createdAt"`
}

// still need to figure out structure of CIK info table
type Info struct {
}

type ReceiveDataPayload struct {
}
