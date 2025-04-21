package service

import (
	"context"

	"1337b04rd/internal/domain/entity"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	GetCommentByID(ctx context.Context, id int) (entity.Comment, error)
	DeleteComment(ctx context.Context, id int) error
}
