package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	FCMToken     string    `json:"fcm_token,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
