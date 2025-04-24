package servicework

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
	"context"
	"fmt"
	"time"
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
	fmt.Println(">>> CreateComment service called")
	comment.CreatedAt = time.Now()

	createdComment, err := s.commentRepo.CreateComment(ctx, &comment)
	if err != nil {
		return entity.Comment{}, err
	}

	return *createdComment, nil
}

func (s *commentService) GetCommentByID(ctx context.Context, id int) ([]entity.Comment, error) {
	fmt.Println(">>> GetCommentById service called", id)
	comments, err := s.commentRepo.GetCommentsByPostID(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, fmt.Errorf("comment not found for postID %d", id)
	}

	var result []entity.Comment
	for _, c := range comments {
		result = append(result, *c)
	}

	return result, nil
}

func (s *commentService) DeleteComment(ctx context.Context, id int) error {
	fmt.Println(">>> DeleteComment service called")
	return s.commentRepo.DeleteComment(ctx, id)
}
