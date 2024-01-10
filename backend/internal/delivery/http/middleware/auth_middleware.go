package middleware

import (
	"backend/internal/delivery/http/exception"
	"backend/internal/model"
	"backend/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func AuthMiddleware(viperConfig *viper.Viper, redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")

			// Get token from Bearer
			if !strings.Contains(authHeader, "Bearer ") {
				return exception.NewUnauthorizedError(exception.EmptyTokenMsg)
			}
			accessToken := strings.Replace(authHeader, "Bearer ", "", -1)

			// Parse access token data and check expiration time
			accessTokenData, err := utils.ParseAccessToken(viperConfig, accessToken)
			if err != nil {
				return err
			}

			// Get token from redis
			accessTokenRedis, err := redisClient.Get(c.Request().Context(), utils.GenerateAccessTokenRedisKey(accessTokenData.UserID)).Result()
			if err != nil {
				return exception.NewUnauthorizedError(exception.InvalidTokenMsg)
			}

			// Verify if the token are same
			if accessToken != accessTokenRedis {
				return exception.NewUnauthorizedError(exception.InvalidTokenMsg)
			}

			// If token are valid, set user data on context
			c.Set("userAuthData", &model.CurrentUser{
				ID:    accessTokenData.UserID,
				Email: accessTokenData.UserEmail,
				Name:  accessTokenData.UserName,
			})

			return next(c)
		}
	}
}
