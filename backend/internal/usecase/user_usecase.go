package usecase

import (
	"backend/internal/delivery/http/exception"
	"backend/internal/entity"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	Config         *viper.Viper
}

func NewUserUseCase(db *gorm.DB, redis *redis.Client, validate *validator.Validate,
	userRepository *repository.UserRepository, config *viper.Viper) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Redis:          redis,
		Validate:       validate,
		UserRepository: userRepository,
		Config:         config,
	}
}

func (s *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return err
	}

	var user entity.User
	if _ = s.UserRepository.FindByEmail(tx, &user, request.Email); len(user.ID) > 0 {
		return exception.NewConflictError("user")
	}

	userPassword, err := utils.HashUserPassword(request.Password)
	if err != nil {
		return err
	}
	user = entity.User{
		ID:        "USR" + utils.GenerateRandomString(20),
		Name:      request.Name,
		Email:     request.Email,
		Password:  userPassword,
		CreatedAt: time.Now(),
	}
	if err := s.UserRepository.Save(tx, &user); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.TokenData, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := s.Validate.Struct(request); err != nil {
		return nil, err
	}

	// Check if user exists
	userFound := new(entity.User)
	if err := s.UserRepository.FindByEmail(tx, userFound, request.Email); err != nil || userFound.ID == "" {
		return nil, exception.NewUnauthorizedError("user not found")
	}

	// If user exits, check password
	if !utils.IsUserPasswordValid(request.Password, userFound.Password) {
		return nil, exception.NewUnauthorizedError("password doesn't match")
	}

	// If password valid, generate token
	response := new(model.TokenData)
	var (
		err           error
		accessExpDur  time.Duration
		refreshExpDur time.Duration
	)

	response.AccessToken, accessExpDur, err = utils.GenerateAccessToken(s.Config, userFound)
	if err != nil {
		return nil, err
	}

	response.RefreshToken, refreshExpDur, err = utils.GenerateRefreshToken(s.Config, userFound)
	if err != nil {
		return nil, err
	}

	response.RefreshExpAt = time.Now().Add(refreshExpDur)

	// Store both in redis
	if err := s.Redis.SetEx(ctx, utils.GenerateAccessTokenRedisKey(userFound.ID),
		response.AccessToken, accessExpDur).Err(); err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	if err := s.Redis.SetEx(ctx, utils.GenerateRefreshTokenRedisKey(userFound.ID),
		response.RefreshToken, refreshExpDur).Err(); err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	return response, nil
}

func (s *UserUseCase) Current(ctx context.Context, currentUser *model.CurrentUser) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := s.UserRepository.Repository.FindByID(tx, user, currentUser.ID); err != nil {
		return err
	}

	currentUser.Name = user.Email
	currentUser.CreatedAt = user.CreatedAt

	return nil
}

func (s *UserUseCase) Refresh(ctx context.Context, tokenData *model.TokenData) error {
	// Check if refresh token is not expired and parse it
	userAuthData, err := utils.ParseRefreshToken(s.Config, tokenData.RefreshToken)
	if err != nil {
		return err
	}

	// Check if refresh is available in redis
	refreshTokenRedisKey := utils.GenerateRefreshTokenRedisKey(userAuthData.UserID)
	redisRefreshToken, err := s.Redis.Get(ctx, refreshTokenRedisKey).Result()
	if err != nil {
		return exception.NewUnauthorizedError(exception.InvalidTokenMsg)
	}

	// Check if request refresh token is same as redis
	if tokenData.RefreshToken != redisRefreshToken {
		return exception.NewUnauthorizedError(exception.InvalidTokenMsg)
	}

	// Generate new access token
	newAccessToken, expTimeAccessTokenDur, err := utils.GenerateAccessToken(s.Config, &entity.User{
		ID:    userAuthData.UserID,
		Email: userAuthData.UserEmail,
	})
	tokenData.AccessToken = newAccessToken
	// Set new access token on redis
	accessTokenRedisKey := utils.GenerateAccessTokenRedisKey(userAuthData.UserID)
	if err := s.Redis.SetEx(ctx, accessTokenRedisKey, newAccessToken, expTimeAccessTokenDur).Err(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *UserUseCase) Logout(ctx context.Context, currentUser *model.CurrentUser) error {
	// Delete refresh token from redis
	refreshTokenRedisKey := utils.GenerateAccessTokenRedisKey(currentUser.ID)
	if err := s.Redis.Del(ctx, refreshTokenRedisKey).Err(); err != nil {
		return err
	}

	// Delete acess token from redis
	accessTokenRedisKey := utils.GenerateRefreshTokenRedisKey(currentUser.ID)
	if err := s.Redis.Del(ctx, accessTokenRedisKey).Err(); err != nil {
		return err
	}

	return nil
}
