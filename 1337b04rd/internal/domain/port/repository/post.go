package repository

import (
	"1337b04rd/internal/domain/entity"
)

type PostRepository interface {
	CreatePost(post *entity.Post) (*entity.Post, error)
	GetPostByID(id int) (*entity.Post, error)
	UpdatePost(post *entity.Post) (*entity.Post, error)
	DeletePost(id int) error
}
