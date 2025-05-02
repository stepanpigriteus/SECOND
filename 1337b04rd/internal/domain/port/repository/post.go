package repository

import (
	"context"

	"a1337b04rd/internal/domain/entity"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	GetPostByID(ctx context.Context, id int32) (*entity.Post, error)
	UpdatePost(ctx context.Context, postID int) error
	// UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id int32) error
	ListPosts(ctx context.Context) ([]entity.PostRequest, error)
	ListArchivedPosts(ctx context.Context) ([]entity.PostRequest, error)
}
