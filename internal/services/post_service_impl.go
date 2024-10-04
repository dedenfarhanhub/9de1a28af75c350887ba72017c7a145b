package services

import (
	"errors"
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/entities"
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"github.com/dedenfarhanhub/blog-service/internal/repositories"
	"time"
)

// PostServiceImpl struct
type PostServiceImpl struct {
	postRepo     repositories.PostRepository
	userService  UserService
	redisService *RedisService
}

// CreatePost creates a new post
func (s *PostServiceImpl) CreatePost(postRequest *dto.PostRequest) (*dto.PostResponse, error) {
	if err := s.validatePostRequest(postRequest); err != nil {
		return nil, err
	}

	author, err := s.userService.FindAuthorByID(postRequest.AuthorID)
	if err != nil {
		return nil, err
	}

	if author == nil {
		return nil, errors.New("author not found")
	}

	postEntity := s.newPostEntity(postRequest, author)

	if err := s.createPost(postEntity); err != nil {
		return nil, err
	}

	if err := s.cachePost(postEntity); err != nil {
		return nil, err
	}

	return postEntity.ToPostResponse(author.ToAuthorResponse()), nil
}

// GetPostByID retrieves a post by its ID
func (s *PostServiceImpl) GetPostByID(id uint) (*dto.PostResponse, error) {
	post, err := s.getPostFromCache(id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		post, err = s.getPostFromDatabase(id)
		if err != nil {
			return nil, err
		}
		if post == nil {
			return nil, errors.New("post not found")
		}

		if err := s.cachePost(post); err != nil {
			return nil, err
		}
	}

	return post.ToPostResponse(post.Author.ToAuthorResponse()), nil
}

// GetAll func
func (s *PostServiceImpl) GetAll(params *dto.QueryParams) ([]*dto.PostResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10 // Default page size
	}

	posts, err := s.postRepo.FindAllWithFilters(params)
	if err != nil {
		return nil, err
	}

	var postResponses []*dto.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToPostResponse(post.Author.ToAuthorResponse()))
	}

	return postResponses, nil
}

// Update func
func (s *PostServiceImpl) Update(id uint, postRequest *dto.PostRequest) (*dto.PostResponse, error) {
	existingPost, err := s.getPostWithOwnershipCheck(id, postRequest.AuthorID)
	if err != nil {
		return nil, err
	}

	existingPost.Title = postRequest.Title
	existingPost.Content = postRequest.Content
	existingPost.UpdatedAt = time.Now()
	if err := s.updatePost(existingPost); err != nil {
		return nil, err
	}
	if err := s.cachePost(existingPost); err != nil {
		return nil, err
	}

	return existingPost.ToPostResponse(existingPost.Author.ToAuthorResponse()), nil
}

// Delete func
func (s *PostServiceImpl) Delete(id uint, userID uint) error {
	existingPost, err := s.getPostWithOwnershipCheck(id, userID)
	if err != nil {
		return err
	}

	// Delete the post from the database
	if err := s.postRepo.Delete(existingPost.ID); err != nil {
		return err
	}

	// Remove the post from Redis cache
	idStr, _ := helpers.ConvertToString(existingPost.ID)
	if err := s.redisService.DeleteEntity("post", idStr); err != nil {
		return err
	}
	return nil
}

// Count func
func (s *PostServiceImpl) Count(params *dto.QueryParams) (int64, error) {
	return s.postRepo.Count(params)
}

// NewPostService initializes post service
func NewPostService(postRepo repositories.PostRepository, userService UserService, redisService *RedisService) PostService {
	return &PostServiceImpl{
		postRepo:     postRepo,
		userService:  userService,
		redisService: redisService,
	}
}

// validatePostRequest validates the post request parameters
func (s *PostServiceImpl) validatePostRequest(postRequest *dto.PostRequest) error {
	if postRequest.Title == "" || postRequest.Content == "" {
		return errors.New("title and content are required")
	}
	return nil
}

// newPostEntity creates a new post entity from the request and author
func (s *PostServiceImpl) newPostEntity(postRequest *dto.PostRequest, author *entities.User) *entities.Post {
	return &entities.Post{
		Title:     postRequest.Title,
		Content:   postRequest.Content,
		AuthorID:  postRequest.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Author:    author,
	}
}

// createPost crate the post to the database
func (s *PostServiceImpl) createPost(postEntity *entities.Post) error {
	if err := s.postRepo.Create(postEntity); err != nil {
		return errors.New("failed to create post")
	}
	return nil
}

// updatePost update the post to the database
func (s *PostServiceImpl) updatePost(postEntity *entities.Post) error {
	if err := s.postRepo.Update(postEntity); err != nil {
		return errors.New("failed to create post")
	}
	return nil
}

// cachePost stores the post in Redis for caching
func (s *PostServiceImpl) cachePost(postEntity *entities.Post) error {
	idStr, _ := helpers.ConvertToString(postEntity.ID)
	if err := s.redisService.SetEntity("post", idStr, postEntity, 24*time.Hour); err != nil {
		return errors.New("failed to store post in Redis")
	}
	return nil
}

// getPostFromCache retrieves a post from Redis
func (s *PostServiceImpl) getPostFromCache(id uint) (*entities.Post, error) {
	idStr, _ := helpers.ConvertToString(id)
	existingPost := &entities.Post{}
	err := s.redisService.GetEntity("post", idStr, existingPost)
	if err != nil {
		return nil, errors.New("failed to retrieve post from Redis")
	}
	if existingPost.ID != 0 {
		if existingPost.Author.ID == 0 {
			// If Author is missing, we need to load it from the database
			author, err := s.userService.FindAuthorByID(existingPost.AuthorID)
			if err != nil {
				return nil, errors.New("failed to load author details")
			}
			existingPost.Author = author
		}
		return existingPost, nil
	}
	return nil, nil
}

// getPostFromDatabase retrieves a post from the database
func (s *PostServiceImpl) getPostFromDatabase(id uint) (*entities.Post, error) {
	postFromDB, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("failed to find post in database")
	}
	return postFromDB, nil
}

// getPostWithOwnershipCheck retrieves a post and checks if the user is the author
func (s *PostServiceImpl) getPostWithOwnershipCheck(id uint, userID uint) (*entities.Post, error) {
	// Check Redis first
	post, err := s.getPostFromCache(id)
	if err != nil {
		return nil, err
	}

	// If not found in Redis, check the database
	if post == nil {
		post, err = s.getPostFromDatabase(id)
		if err != nil {
			return nil, err
		}
		if post == nil {
			return nil, errors.New("post not found")
		}
	}

	// Check ownership
	if post.AuthorID != userID {
		return nil, errors.New("you do not have permission to modify this post")
	}

	return post, nil
}
