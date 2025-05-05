package customerman

import "github.com/lonewolf101101/Architect-betting/backend/pkg/entities"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type Customer struct {
	entities.Model
	GoogleID       string `json:"google_id"`
	ProfilePicture string `json:"profile_picture"`
	Email          string `json:"email" gorm:"index:idx_user_email"`
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone_number"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}
