package entities

import "time"

type Click struct {
	ID        int       `json:"id"`
	LinkID    int       `json:"link_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
