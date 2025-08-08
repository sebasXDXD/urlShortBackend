package entities

import "time"

type LinkStats struct {
	TotalCreated     int
	TotalClicks      int
	MostClickedURL   string
	MostClickedCount int
	LastAccess       time.Time
}
