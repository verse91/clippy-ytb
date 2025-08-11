package model

type UserProfile struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Credits int    `json:"credits"`
}
