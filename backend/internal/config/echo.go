package config

import (
	"backend/internal/delivery/http/exception"

	"github.com/labstack/echo/v4"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = exception.CustomErrorHandler

	return e
}
