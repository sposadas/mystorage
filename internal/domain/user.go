package domain

import "github.com/google/uuid"

type User struct {
	ID       int       `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Active   bool      `json:"active"`
}
