package test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

func newRequest(method string, url string, requestBody string) *http.Request {
	request := httptest.NewRequest(method, url, strings.NewReader(requestBody))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return request
}

func newRequestWithToken(method string, url string, requestBody string, token string) *http.Request {
	bearer := "Bearer " + token
	request := newRequest(method, url, requestBody)
	request.Header.Add("Authorization", bearer)

	return request
}
