package entity

import "time"

type User struct {
	ID            int       `json:"id"`
	AvatarURL     string    `json:"avatar_url"`
	CharacterName string    `json:"character_name"`
	CustomName    string    `json:"custom_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// {
// 	"id": 1,
// 	"avatar_url": "https://example.com/avatar.jpg",
// 	"character_name": "Hero123",
// 	"custom_name": "SuperHeroX",
// 	"created_at": "2025-04-24T08:00:00Z"
//   }
