package repositories

import (
	"errors"
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
	"gorm.io/gorm"
)

// PostRepository interface
type PostRepository interface {
	Create(post *entities.Post) error
	FindByID(id uint) (*entities.Post, error)
	FindAll() ([]entities.Post, error)
	Update(post *entities.Post) error
	Delete(id uint) error
	FindAllWithFilters(params *dto.QueryParams) ([]entities.Post, error)
	Count(params *dto.QueryParams) (int64, error)
}

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository initialize post repository
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *entities.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*entities.Post, error) {
	var post entities.Post
	if err := r.db.Preload("Author").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Mengembalikan nil jika user tidak ditemukan
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll() ([]entities.Post, error) {
	var posts []entities.Post
	err := r.db.Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(post *entities.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Post{}, id).Error
}

func (r *postRepository) FindAllWithFilters(params *dto.QueryParams) ([]entities.Post, error) {
	var posts []entities.Post
	query := r.db.Model(&entities.Post{}).Preload("Author")

	// Apply search filter
	if params.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) Count(params *dto.QueryParams) (int64, error) {
	var count int64
	query := r.db.Model(&entities.Post{})

	// Apply search filter
	if params.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Count the total
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
