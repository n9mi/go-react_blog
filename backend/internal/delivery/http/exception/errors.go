package exception

import (
	"github.com/labstack/echo/v4"
)

func NewNotFoundError(entity string) error {
	err := echo.ErrNotFound
	err.Message = entity + " is not found"

	return err
}

func NewBadRequestError(message string) error {
	err := echo.ErrBadRequest
	err.Message = message

	return err
}

func NewUnauthorizedError(message string) error {
	err := echo.ErrUnauthorized
	err.Message = message

	return err
}

func NewForbiddenError(message string) error {
	err := echo.ErrForbidden
	err.Message = message

	return err
}

func NewConflictError(entity string) error {
	err := echo.ErrConflict
	err.Message = entity + " is alredy exists"

	return err
}

func NewInternalServerError(message string) error {
	err := echo.ErrInternalServerError
	err.Message = message

	return err
}
