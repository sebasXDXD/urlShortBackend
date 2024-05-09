package entities

import (
	"time"
)

// Link representa un enlace en la aplicación.
type Link struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	RedirectTo    string    `json:"redirect_to"`
	UserCreatedID int       `json:"user_created_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsDeleted     bool      `json:"is_deleted"`
	DeletedAt     time.Time `json:"deleted_at"`
}
