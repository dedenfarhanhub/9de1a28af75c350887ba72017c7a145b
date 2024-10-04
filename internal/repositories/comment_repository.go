package repositories

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
	"gorm.io/gorm"
)

// CommentRepository interface
type CommentRepository interface {
	Create(comment *entities.Comment) error
	FindAllByPostIDWithFilters(postID uint, params *dto.QueryParams) ([]entities.Comment, error)
	CountByPostID(postID uint, params *dto.QueryParams) (int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository initializes comment repository
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *entities.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindAllByPostIDWithFilters(postID uint, params *dto.QueryParams) ([]entities.Comment, error) {
	var comments []entities.Comment
	query := r.db.Model(&entities.Comment{}).Where("post_id = ?", postID)

	// Apply search filter
	if params.Search != "" {
		query = query.Where("author_name LIKE ? OR content LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Apply sorting
	if params.SortBy != "" {
		if params.SortOrder == "desc" {
			query = query.Order(params.SortBy + " desc")
		} else {
			query = query.Order(params.SortBy + " asc")
		}
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) CountByPostID(postID uint, params *dto.QueryParams) (int64, error) {
	var count int64
	query := r.db.Model(&entities.Comment{}).Where("post_id = ?", postID)

	// Apply search filter
	if params.Search != "" {
		query = query.Where("author_name LIKE ? OR content LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Count the total
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
