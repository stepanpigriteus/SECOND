package servicework

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

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

func (s *commentService) CreateComment(ctx context.Context, comment entity.Comment, fileHeader *multipart.FileHeader) (entity.Comment, error) {
	fmt.Println(">>> CreateComment service called")

	comment.CreatedAt = time.Now()

	// Если файл есть, сохранить его через порт (интерфейс)
	// if file != nil {
	// 	fileURL, err := s.fileStorage.SaveFile(ctx, file, fileHeader.Filename)
	// 	if err != nil {
	// 		return entity.Comment{}, err
	// 	}
	// 	comment.FileURL = fileURL
	// }

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
