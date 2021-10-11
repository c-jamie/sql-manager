package api

import (
	"time"
)

type UserAccount struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Activated   bool      `json:"activated"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int       `json:"-"`
	Team        *Team     `json:"team"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
}
