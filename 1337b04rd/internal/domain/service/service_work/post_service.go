package servicework

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
	"1337b04rd/internal/domain/service"
	"context"
	"fmt"
	"time"
)

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) service.PostService {
	return &postService{postRepo: repo}
}

func (s *postService) CreatePost(ctx context.Context, post entity.Post) (entity.Post, error) {
	fmt.Println(">>> CreatePost service called")
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	createdPost, err := s.postRepo.CreatePost(ctx, &post)
	if err != nil {
		return entity.Post{}, err
	}

	return *createdPost, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (entity.Post, error) {
	fmt.Println(">>> GetPostById service called", id)
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return entity.Post{}, err
	}
	return *post, nil
}

func (s *postService) UpdatePost(ctx context.Context, id int, post entity.Post) (entity.Post, error) {
	fmt.Println(">>> UpdatePost service called")
	post.ID = id
	post.UpdatedAt = time.Now()

	updatedPost, err := s.postRepo.UpdatePost(ctx, &post)
	if err != nil {
		return entity.Post{}, err
	}
	return *updatedPost, nil
}

func (s *postService) DeletePost(ctx context.Context, id int) error {
	fmt.Println(">>> DeletePost service called")
	return s.postRepo.DeletePost(ctx, id)
}
