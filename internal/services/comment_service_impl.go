package services

import (
	"errors"
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
	"github.com/dedenfarhanhub/blog-service/internal/repositories"
)

// CommentServiceImpl struct
type CommentServiceImpl struct {
	commentRepo repositories.CommentRepository
	postService PostService
}

// Create func create comment
func (s *CommentServiceImpl) Create(postID uint, commentRequest *dto.CommentRequest) (*dto.CommentResponse, error) {
	// Validate if the PostID exists
	post, err := s.postService.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("postID does not exist")
	}

	comment := &entities.Comment{
		PostID:     postID,
		AuthorName: commentRequest.AuthorName,
		Content:    commentRequest.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return comment.ToCommentResponse(), nil
}

// GetAllByPostID get all comments by post id
func (s *CommentServiceImpl) GetAllByPostID(postID uint, params *dto.QueryParams) ([]*dto.CommentResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10 // Default page size
	}

	comments, err := s.commentRepo.FindAllByPostIDWithFilters(postID, params)
	if err != nil {
		return nil, err
	}

	var commentResponses []*dto.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, comment.ToCommentResponse())
	}

	return commentResponses, nil
}

// CountAllByPostID count all comments by post id
func (s *CommentServiceImpl) CountAllByPostID(postID uint, params *dto.QueryParams) (int64, error) {
	return s.commentRepo.CountByPostID(postID, params)
}

// NewCommentService initializes comment service
func NewCommentService(commentRepo repositories.CommentRepository, postService PostService) CommentService {
	return &CommentServiceImpl{
		commentRepo: commentRepo,
		postService: postService,
	}
}
