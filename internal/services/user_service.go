package services

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
)

// UserService interface
type UserService interface {
	Register(userRequest *dto.UserRequest) (*dto.UserResponse, error)
	HashPassword(password string) (string, error)
	Login(userLoginRequest *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	FindAuthorByID(authorID uint) (*entities.User, error)
}
