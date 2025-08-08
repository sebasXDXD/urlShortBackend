package dtos

import "time"

type UserStats struct {
	TotalCreated     int
	TotalClicks      int
	MostClickedURL   string
	MostClickedCount int
	LastAccess       time.Time
}
