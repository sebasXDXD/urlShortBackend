package entities

import "time"

type UserProfileResponse struct {
	FullName string    `json:"fullName"`
	Email    string    `json:"email"`
	Company  string    `json:"company"`
	Country  string    `json:"country"`
	Phone    string    `json:"phone"`
	Plan     string    `json:"plan"`
	JoinDate time.Time `json:"joinDate"`
	Stats    StatsData `json:"stats"`
}

type StatsData struct {
	URLsCreated      int    `json:"urlsCreated"`
	URLsLimit        int    `json:"urlsLimit"`
	TotalClicks      int    `json:"totalClicks"`
	PopularURL       string `json:"popularUrl"`
	PopularURLClicks int    `json:"popularUrlClicks"`
	LastAccess       string `json:"lastAccess"`
}
