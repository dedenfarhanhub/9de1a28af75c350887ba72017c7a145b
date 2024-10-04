package services

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
)

// PostService interface
type PostService interface {
	CreatePost(postRequest *dto.PostRequest) (*dto.PostResponse, error)
	GetPostByID(id uint) (*dto.PostResponse, error)
	GetAll(params *dto.QueryParams) ([]*dto.PostResponse, error)
	Update(id uint, postRequest *dto.PostRequest) (*dto.PostResponse, error)
	Delete(id uint, userID uint) error
	Count(params *dto.QueryParams) (int64, error)
}
