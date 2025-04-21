package servicework

import (
	"context"
	"time"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
)

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) *postService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, post entity.Post) (entity.Post, error) {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	createdPost, err := s.postRepo.CreatePost(ctx, &post)
	if err != nil {
		return entity.Post{}, err
	}

	return *createdPost, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (entity.Post, error) {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return entity.Post{}, err
	}
	return *post, nil
}

func (s *postService) UpdatePost(ctx context.Context, id int, post entity.Post) (entity.Post, error) {
	post.ID = id
	post.UpdatedAt = time.Now()

	updatedPost, err := s.postRepo.UpdatePost(ctx, &post)
	if err != nil {
		return entity.Post{}, err
	}
	return *updatedPost, nil
}

func (s *postService) DeletePost(ctx context.Context, id int) error {
	return s.postRepo.DeletePost(ctx, id)
}
