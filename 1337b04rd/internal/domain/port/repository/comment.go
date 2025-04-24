package repository

import (
	"1337b04rd/internal/domain/entity"
	"context"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, id int) error
}
