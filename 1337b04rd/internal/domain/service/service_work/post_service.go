package servicework

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
	"context"
)

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) *postService {
	return &postService{
		postRepo: postRepo,
	}
}

func CreatePost(ctx context.Context, post entity.Post) (entity.Post, error) {
	return post, nil
}

func GetPostByID(ctx context.Context, id int) (entity.Post, error) {
	var post entity.Post
	return post, nil
}

func UpdatePost(ctx context.Context, id int, post entity.Post) (entity.Post, error) {
	return post, nil
}

func DeletePost(ctx context.Context, id int) error {
	return nil
}
