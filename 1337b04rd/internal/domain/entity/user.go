package entity

import "time"

type User struct {
	ID            int       `json:"id"`
	AvatarURL     string    `json:"avatar_url"`
	CharacterName string    `json:"character_name"`
	CustomName    string    `json:"custom_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
