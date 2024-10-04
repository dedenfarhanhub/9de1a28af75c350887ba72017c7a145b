package repositories

import (
	"errors"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
	"gorm.io/gorm"
)

// UserRepository interface
type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uint) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository initialize user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Mengembalikan nil jika user tidak ditemukan
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Mengembalikan nil jika user tidak ditemukan
		}
		return nil, err
	}
	return &user, nil
}
