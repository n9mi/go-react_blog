package http

import (
	"backend/internal/constant"
	"backend/internal/delivery/http/exception"
	"backend/internal/model"
	"backend/internal/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
}

func NewUserController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{UserUseCase: userUseCase}
}

func (ct *UserController) Register(c echo.Context) error {
	request := new(model.RegisterUserRequest)

	if err := c.Bind(request); err != nil {
		return err
	}

	if err := ct.UserUseCase.Create(c.Request().Context(), request); err != nil {
		return err
	}

	response := model.DataResponse[any]{
		Code:   http.StatusOK,
		Status: "OK",
	}
	return c.JSON(response.Code, response)
}

func (ct *UserController) Login(c echo.Context) error {
	request := new(model.LoginUserRequest)

	if err := c.Bind(request); err != nil {
		return err
	}

	loginResponse, err := ct.UserUseCase.Login(c.Request().Context(), request)
	if err != nil {
		return err
	}

	fmt.Println("REFRESH_EXP_AT", loginResponse.RefreshExpAt)

	// Send refresh token to cookie
	cookie := new(http.Cookie)
	cookie.Name = constant.REFRESH_TOKEN_COOKIE_NAME
	cookie.Value = loginResponse.RefreshToken
	cookie.Expires = loginResponse.RefreshExpAt
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	response := model.DataResponse[*model.TokenData]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   loginResponse,
	}
	return c.JSON(response.Code, response)
}

func (ct *UserController) Current(c echo.Context) error {
	// Get current user ID and email from context
	authDataInterface := c.Get(constant.USER_AUTH_DATA_CONTEXT_NAME)

	// Check if data is on CurrentUser type
	currentUser, ok := authDataInterface.(*model.CurrentUser)
	if !ok {
		return exception.NewUnauthorizedError("")
	}

	// Get full user data from service
	ct.UserUseCase.Current(c.Request().Context(), currentUser)

	response := model.DataResponse[*model.CurrentUser]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   currentUser,
	}
	return c.JSON(response.Code, response)
}

func (ct *UserController) Refresh(c echo.Context) error {
	// Get refresh token from HTTPOnly cookie
	cookie, err := c.Cookie(constant.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		return exception.NewUnauthorizedError(err.Error())
	}

	// Check if token empty
	if cookie.Value == "" {
		return exception.NewUnauthorizedError(exception.EmptyTokenMsg)
	}

	tokenData := new(model.TokenData)
	tokenData.RefreshToken = cookie.Value

	// Get new access token by passing current refresh token
	if err := ct.UserUseCase.Refresh(c.Request().Context(), tokenData); err != nil {
		return err
	}

	response := model.DataResponse[*model.TokenData]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tokenData,
	}
	return c.JSON(response.Code, response)
}

func (ct *UserController) Logout(c echo.Context) error {
	// Get refresh token from HTTPOnly cookie
	cookie, err := c.Cookie(constant.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		return exception.NewUnauthorizedError(err.Error())
	}

	// Delete cookie by setting maxAge to negative value
	cookie.MaxAge = -1

	// Send the cookie
	c.SetCookie(cookie)

	// Get auth data from context
	authDataInterface := c.Get(constant.USER_AUTH_DATA_CONTEXT_NAME)

	currentUser, ok := authDataInterface.(*model.CurrentUser)
	if !ok {
		return exception.NewUnauthorizedError("")
	}

	// Passing auth data to service. Service will remove both of access token and refresh token
	err = ct.UserUseCase.Logout(c.Request().Context(), currentUser)
	if err != nil {
		return err
	}

	return nil
}
