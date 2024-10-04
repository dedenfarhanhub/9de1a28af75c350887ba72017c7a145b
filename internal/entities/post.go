package entities

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"time"
)

// Post represents a blog post.
type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	AuthorID  uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Author   *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments []*Comment `gorm:"foreignKey:PostID"`
}

// ToPostResponse converts a Post entity to a PostResponse DTO.
func (p *Post) ToPostResponse(author *dto.AuthorResponse) *dto.PostResponse {
	return &dto.PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		Author:    author,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
}
