package service

import (
	"1337b04rd/internal/domain/entity"
	"context"
)

type PostService interface {
	CreatePost(ctx context.Context, post entity.Post) (entity.Post, error)
	GetPostByID(ctx context.Context, id int) (entity.Post, error)
	UpdatePost(ctx context.Context, id int, post entity.Post) (entity.Post, error)
	DeletePost(ctx context.Context, id int) error
}
