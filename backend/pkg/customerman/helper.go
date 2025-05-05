package customerman

import "time"

type ToduUser struct {
	Email               string    `json:"email"`
	EmailVerified       bool      `json:"email_verified"`
	Name                string    `json:"name"`
	PhoneNumber         string    `json:"phone_number"`
	PhoneNumberVerified bool      `json:"phone_number_verified"`
	Sub                 string    `json:"sub"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type Filter struct {
	Keyword     string
	CustomerIDs []int
}
