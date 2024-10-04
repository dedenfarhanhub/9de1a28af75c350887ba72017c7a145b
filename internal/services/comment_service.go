package services

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
)

// CommentService interface
type CommentService interface {
	Create(postID uint, commentRequest *dto.CommentRequest) (*dto.CommentResponse, error)
	GetAllByPostID(postID uint, params *dto.QueryParams) ([]*dto.CommentResponse, error)
	CountAllByPostID(postID uint, params *dto.QueryParams) (int64, error)
}
