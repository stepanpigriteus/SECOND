package servicework

import (
	"context"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
)

type commentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) *commentService {
	return &commentService{
		commentRepo: commentRepo,
	}
}

func (s *commentService) CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	return comment, nil
}

func (s *commentService) GetCommentByID(ctx context.Context, id int) (entity.Comment, error) {
	var comment entity.Comment
	return comment, nil
}

func (s *commentService) DeleteComment(ctx context.Context, id int) error {
	return nil
}
