package entities

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"time"
)

// User struct model
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Posts []*Post `gorm:"foreignKey:AuthorID"`
}

// ToUserResponse convert User to UserResponse
func (u *User) ToUserResponse(token string) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Token: token,
	}
}

// ToUserLoginResponse convert User to UserLoginResponse
func (u *User) ToUserLoginResponse(token string) *dto.UserLoginResponse {
	return &dto.UserLoginResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Token: token,
	}
}

// ToAuthorResponse convert User to AuthorResponse
func (u *User) ToAuthorResponse() *dto.AuthorResponse {
	return &dto.AuthorResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
