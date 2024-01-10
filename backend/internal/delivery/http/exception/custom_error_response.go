package exception

import (
	"backend/internal/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func SplitErrorMessage(errorMessage string) string {
	return strings.Split(strings.Split(errorMessage, ",")[1], "=")[1]
}

func GetBadRequestErrorResponse(err error) model.MessagesResponse {
	return model.MessagesResponse{
		Code:     http.StatusBadRequest,
		Status:   "BAD REQUEST",
		Messages: []string{SplitErrorMessage(err.Error())},
	}
}

func GetNotFoundErrorResponse(err error) model.MessagesResponse {
	return model.MessagesResponse{
		Code:     http.StatusNotFound,
		Status:   "NOT FOUND",
		Messages: []string{SplitErrorMessage(err.Error())},
	}
}

func GetValidationErrorResponse(errConv validator.ValidationErrors) model.MessagesResponse {
	var response = model.MessagesResponse{
		Code:   http.StatusBadRequest,
		Status: "BAD REQUEST",
	}

	for _, errItem := range errConv {
		switch errItem.Tag() {
		case "required":
			response.Messages = append(response.Messages,
				fmt.Sprintf("%s is required", errItem.Field()))
		case "min":
			response.Messages = append(response.Messages,
				fmt.Sprintf("%s is should be more than %s character",
					errItem.Field(), errItem.Param()))
		case "max":
			response.Messages = append(response.Messages,
				fmt.Sprintf("%s is should be less than %s character",
					errItem.Field(), errItem.Param()))
		case "email":
			response.Messages = append(response.Messages,
				fmt.Sprintf("%s should be a valid email", errItem.Field()))
		}
	}

	return response
}

func GetUnauthorizedErrorResponse(err error) model.MessagesResponse {
	return model.MessagesResponse{
		Code:     http.StatusUnauthorized,
		Status:   "UNAUTHORIZED",
		Messages: []string{SplitErrorMessage(err.Error())},
	}
}

func GetForbiddenErrorResponse(err error) model.MessagesResponse {
	return model.MessagesResponse{
		Code:     http.StatusForbidden,
		Status:   "FORBIDDEN",
		Messages: []string{SplitErrorMessage(err.Error())},
	}
}

func GetConflictErrorResponse(err error) model.MessagesResponse {
	return model.MessagesResponse{
		Code:     http.StatusConflict,
		Status:   "CONFLICT",
		Messages: []string{SplitErrorMessage(err.Error())},
	}
}

func GetInternalServerError(err error) model.MessagesResponse {
	return model.MessagesResponse{
		Code:     http.StatusInternalServerError,
		Status:   "INTERNAL SERVER ERROR",
		Messages: []string{SplitErrorMessage(err.Error())},
	}
}
