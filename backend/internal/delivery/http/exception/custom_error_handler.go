package exception

import (
	"backend/internal/model"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func CustomErrorHandler(err error, c echo.Context) {
	var response model.MessagesResponse

	if errors.Is(err, echo.ErrNotFound) {
		response = GetNotFoundErrorResponse(err)
	} else if errors.Is(err, echo.ErrBadRequest) {
		response = GetBadRequestErrorResponse(err)
	} else if errConv, ok := err.(validator.ValidationErrors); ok {
		response = GetValidationErrorResponse(errConv)
	} else if errors.Is(err, echo.ErrUnauthorized) {
		response = GetUnauthorizedErrorResponse(err)
	} else if errors.Is(err, echo.ErrForbidden) {
		response = GetForbiddenErrorResponse(err)
	} else if errors.Is(err, echo.ErrConflict) {
		response = GetConflictErrorResponse(err)
	} else if errors.Is(err, echo.ErrInternalServerError) {
		response = GetInternalServerError(err)
	} else {
		if he, ok := err.(*echo.HTTPError); ok {
			switch {
			case he.Code == http.StatusBadRequest:
				response = GetBadRequestErrorResponse(he)
			case he.Code == http.StatusUnauthorized:
				response = GetUnauthorizedErrorResponse(he)
			case he.Code == http.StatusForbidden:
				response = GetForbiddenErrorResponse(he)
			case he.Code == http.StatusBadRequest:
				response = GetBadRequestErrorResponse(he)
			case he.Code == http.StatusNotFound:
				response = GetNotFoundErrorResponse(he)
			case he.Code == http.StatusInternalServerError:
				response = GetInternalServerError(he)
			}
		} else {
			response = GetInternalServerError(err)
		}
	}

	c.JSON(response.Code, response)
}
