package utils

import (
	"backend/internal/delivery/http/exception"
	"backend/internal/entity"
	"backend/internal/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func GenerateAccessToken(viperConfig *viper.Viper, user *entity.User) (string, time.Duration, error) {
	key := viperConfig.GetString("auth.accessTokenKey")
	expMinutes := viperConfig.GetInt("auth.accessTokenExpMinutes")

	return GenerateAuthToken(user, key, expMinutes)
}

func GenerateRefreshToken(viperConfig *viper.Viper, user *entity.User) (string, time.Duration, error) {
	key := viperConfig.GetString("auth.refreshTokenKey")
	expMinutes := viperConfig.GetInt("auth.refreshTokenExpMinutes")

	return GenerateAuthToken(user, key, expMinutes)
}

func GenerateAuthToken(user *entity.User, key string, expMinutes int) (string, time.Duration, error) {
	timeDuration := time.Duration(expMinutes) * time.Minute
	timeExp := time.Now().Add(timeDuration)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": timeExp.Unix(),
		"data": model.UserAuthData{
			UserID:    user.ID,
			UserEmail: user.Email,
			UserName:  user.Name,
		},
	})

	s, err := t.SignedString([]byte(key))

	return s, timeDuration, err
}

func ParseAccessToken(viperConfig *viper.Viper, token string) (*model.UserAuthData, error) {
	key := viperConfig.GetString("auth.accessTokenKey")

	return ParseAuthToken(token, key)
}

func ParseRefreshToken(viperConfig *viper.Viper, token string) (*model.UserAuthData, error) {
	key := viperConfig.GetString("auth.refreshTokenKey")

	return ParseAuthToken(token, key)
}

func ParseAuthToken(token string, key string) (*model.UserAuthData, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, exception.NewUnauthorizedError(exception.InvalidTokenMsg)
	}

	claims := t.Claims.(jwt.MapClaims)
	timeExp := int(claims["exp"].(float64))
	dataInterface := claims["data"].(map[string]interface{})

	// Check if token expired
	if time.Now().Unix() > int64(timeExp) {
		return nil, exception.NewUnauthorizedError(exception.ExpiredTokenMsg)
	}

	// Parse UserAuthData one by one
	userAuthData := new(model.UserAuthData)
	var ok bool

	userAuthData.UserID, ok = dataInterface["UserID"].(string)
	if !ok {
		return nil, exception.NewUnauthorizedError(exception.InvalidTokenMsg)
	}

	userAuthData.UserEmail, ok = dataInterface["UserEmail"].(string)
	if !ok {
		return nil, exception.NewUnauthorizedError(exception.InvalidTokenMsg)
	}

	userAuthData.UserName, ok = dataInterface["UserName"].(string)
	if !ok {
		return nil, exception.NewUnauthorizedError(exception.InvalidTokenMsg)
	}

	return userAuthData, nil
}

func GenerateAccessTokenRedisKey(userID string) string {
	return fmt.Sprintf("ACCESS:%s", userID)
}

func GenerateRefreshTokenRedisKey(userID string) string {
	return fmt.Sprintf("REFRESH:%s", userID)
}
