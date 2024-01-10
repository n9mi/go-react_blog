package usecase

import (
	"backend/internal/delivery/http/exception"
	"backend/internal/entity"
	"backend/internal/model"
	"backend/internal/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostUseCase struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Validate       *validator.Validate
	PostRepository *repository.PostRepository
	UserRepository *repository.UserRepository
}

func NewPostUseCase(db *gorm.DB, redis *redis.Client, validate *validator.Validate, postRepository *repository.PostRepository,
	userRepository *repository.UserRepository) *PostUseCase {
	return &PostUseCase{
		DB:             db,
		Redis:          redis,
		Validate:       validate,
		PostRepository: postRepository,
		UserRepository: userRepository,
	}
}

func (s *PostUseCase) List(ctx context.Context, request *model.PostListRequest) ([]model.PostResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	response, err := s.PostRepository.List(tx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PostUseCase) GetByID(ctx context.Context, request *model.PostGetByIDRequest) (*model.PostResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, err
	}

	response := new(model.PostResponse)
	if err := s.PostRepository.GetWithAuthor(tx, response, request.ID); err != nil {
		return nil, err
	}
	if response.ID == 0 {
		return nil, exception.NewNotFoundError("post")
	}

	return response, nil
}

func (s *PostUseCase) Create(ctx context.Context, request *model.PostCreateRequest) (*model.PostResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()

	// Validate request
	if err := s.Validate.Struct(request); err != nil {
		return nil, err
	}

	// Make entity from request
	post := new(entity.Post)
	post.Title = request.Title
	post.Content = request.Content
	post.UserID = request.AuthorID

	// Save post with repository
	if err := s.PostRepository.Repository.Save(tx, post); err != nil {
		if err := tx.Rollback().Error; err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Confirm created post by retrieving created post from ID
	tx = s.DB.WithContext(ctx).Begin()

	response := new(model.PostResponse)
	if err := s.PostRepository.GetWithAuthor(tx, response, post.ID); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PostUseCase) Update(ctx context.Context, request *model.PostUpdateRequest) (*model.PostResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()

	// Validate request
	if err := s.Validate.Struct(request); err != nil {
		return nil, err
	}

	// Check if post exists, by confirming is this post has same author with current user
	post := new(entity.Post)
	if err := s.PostRepository.GetByIDandAuthorID(tx, post, request.ID, request.AuthorID); err != nil {
		return nil, exception.NewNotFoundError("post")
	}

	// Make entity from request
	post.Title = request.Title
	post.Content = request.Content

	// Save post with repository
	if err := s.PostRepository.Repository.Save(tx, post); err != nil {
		if err := tx.Rollback().Error; err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Confirm updated post by retrieving updated post from ID
	tx = s.DB.WithContext(ctx).Begin()

	response := new(model.PostResponse)
	if err := s.PostRepository.GetWithAuthor(tx, response, post.ID); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PostUseCase) Delete(ctx context.Context, request *model.PostDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := s.Validate.Struct(request); err != nil {
		return nil
	}

	// Check if post exists, by confirming is this post has same author with current user
	post := new(entity.Post)
	if err := s.PostRepository.GetByIDandAuthorID(tx, post, request.ID, request.UserID); err != nil {
		return exception.NewNotFoundError("post")
	}

	// If post exists, delete the post
	if err := s.PostRepository.Repository.Delete(tx, post); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
