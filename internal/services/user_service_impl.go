package services

import (
	"errors"
	"github.com/dedenfarhanhub/blog-service/config"
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"github.com/dedenfarhanhub/blog-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// UserServiceImpl struct
type UserServiceImpl struct {
	userRepo     repositories.UserRepository
	redisService *RedisService
}

// NewUserService initialize user service
func NewUserService(userRepo repositories.UserRepository, redisService *RedisService) UserService {
	return &UserServiceImpl{
		userRepo:     userRepo,
		redisService: redisService,
	}
}

// Register a new user
func (s *UserServiceImpl) Register(userRequest *dto.UserRequest) (*dto.UserResponse, error) {
	// Validate email and password
	if err := validateEmail(userRequest.Email); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(userRequest.Password)
	if err != nil {
		return nil, errors.New("password hashing failed")
	}

	// Check if user exists in Redis
	existingUser, err := s.findUserByEmail(userRequest.Email)
	if err != nil {
		return nil, errors.New("failed to check existing email")
	}
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Create a new user entity
	userEntity := &entities.User{
		Name:         userRequest.Name,
		Email:        userRequest.Email,
		PasswordHash: hashedPassword,
	}

	// Save the new user to the database
	if err := s.userRepo.Create(userEntity); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Store the new user in Redis
	if err := s.redisService.SetEntity("user", userEntity.Email, userEntity, 24*time.Hour); err != nil {
		return nil, errors.New("failed to store user in Redis")
	}

	// Generate JWT token
	cfg := config.LoadConfig()
	token, err := helpers.GenerateToken(userEntity.Email, userEntity.ID, []byte(cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	return userEntity.ToUserResponse(token), nil
}

// HashPassword hashes the user's password
func (s *UserServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Login authenticates the user and returns the user entity
func (s *UserServiceImpl) Login(userLoginRequest *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	// Validate email
	if err := validateEmail(userLoginRequest.Email); err != nil {
		return nil, err
	}

	// Check Redis for existing user
	existingUser, err := s.findUserByEmail(userLoginRequest.Email)
	if err != nil {
		return nil, errors.New("failed to check existing email")
	}
	if existingUser == nil {
		return nil, errors.New("invalid credentials")
	}

	// Validate the password
	if bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(userLoginRequest.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	cfg := config.LoadConfig()
	token, err := helpers.GenerateToken(existingUser.Email, existingUser.ID, []byte(cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	return existingUser.ToUserLoginResponse(token), nil
}

// FindAuthorByID fetches the author (user) by their ID
func (s *UserServiceImpl) FindAuthorByID(authorID uint) (*entities.User, error) {
	author := &entities.User{}
	authorIDStr, _ := helpers.ConvertToString(authorID)
	// Attempt to get the user (author) from Redis first
	err := s.redisService.GetEntity("user", authorIDStr, author)
	if err == nil && author.ID != 0 {
		return author, nil
	}

	// If not found in Redis, check the database
	author, err = s.userRepo.FindByID(authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	// Optionally cache the author in Redis for future use
	_ = s.redisService.SetEntity("user", authorIDStr, author, 24*time.Hour)

	return author, nil
}

// findUserByEmail checks Redis and database for an existing user
func (s *UserServiceImpl) findUserByEmail(email string) (*entities.User, error) {
	existingUser := &entities.User{}
	err := s.redisService.GetEntity("user", email, existingUser)
	if err != nil {
		return nil, errors.New("failed to check existing email in Redis")
	}

	if existingUser.Email == "" {
		// If not found in Redis, check in the database
		existingUserFromDB, err := s.userRepo.FindByEmail(email)
		if err != nil {
			return nil, errors.New("failed to check existing email in database")
		}
		if existingUserFromDB != nil {
			// Store user in Redis for future requests
			err = s.redisService.SetEntity("user", existingUserFromDB.Email, existingUserFromDB, 24*time.Hour)
			if err != nil {
				return existingUserFromDB, err
			}
		}
		return nil, nil
	}
	return existingUser, nil
}

// validateEmail checks if the email format is valid
func validateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil || !matched {
		return errors.New("invalid email format")
	}
	return nil
}
