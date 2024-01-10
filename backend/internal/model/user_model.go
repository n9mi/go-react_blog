package model

import "time"

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=5"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"-"`
	RefreshExpAt time.Time `json:"-"`
}

type UserAuthData struct {
	UserID    string
	UserEmail string
	UserName  string
}

type CurrentUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
