package repository

import "1337b04rd/internal/domain/entity"

type CommentRepository interface {
	CreateComment(comment *entity.Comment) (*entity.Comment, error)
	GetCommentsByPostID(postID int) ([]*entity.Comment, error)
	DeleteComment(id int) error
}
