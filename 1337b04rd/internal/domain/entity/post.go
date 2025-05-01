package entity

import "time"

// Структура для поста
type Post struct {
	ID            int32     `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	ImageURL      string    `json:"image_url,omitempty"`
	UserID        int       `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at,omitempty"`
	LastCommentAt time.Time `json:"last_comment_at"`
}

type PostRequest struct {
	ID       int32  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty"`
}

// {
// 	"title": "My First Post",
// 	"content": "This is the content of my first post.",
// 	"image_url": "https://example.com/image.jpg",
// 	"user_id": 1,
// 	"last_comment_at": "2025-04-24T14:30:00Z"
//   }
