package entities

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"time"
)

// Comment represents a comment on a blog post.
type Comment struct {
	ID         uint      `gorm:"primaryKey"`
	PostID     uint      `gorm:"not null"`
	AuthorName string    `gorm:"not null"`
	Content    string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Post *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// ToCommentResponse converts a Comment entity to a CommentResponse DTO.
func (c *Comment) ToCommentResponse() *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:         c.ID,
		PostID:     c.PostID,
		AuthorName: c.AuthorName,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
	}
}
