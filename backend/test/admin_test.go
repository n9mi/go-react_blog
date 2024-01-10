package test

import (
	"backend/internal/model"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	registerUrl = "http://127.0.0.1:5000/api/auth/register"
	loginUrl    = "http://127.0.0.1:5000/api/auth/login"
)

func TestRegister(t *testing.T) {
	testItems := map[string]TestSchema{
		"USER_Register_OK": {
			"request_name":     "John Doe",
			"request_email":    "johndoe@mail.com",
			"request_password": "johndoe",
			"expected_code":    http.StatusOK,
			"expected_status":  "OK",
		},
		"USER_Register_DUPLICATE": {
			"request_name":     "John Doe",
			"request_email":    "johndoe@mail.com",
			"request_password": "johndoe",
			"expected_code":    http.StatusConflict,
			"expected_status":  "CONFLICT",
		},
		"USER_Register_VALIDATION_ERROR_name_empty": {
			"request_name":     "",
			"request_email":    "johndoe@mail.com",
			"request_password": "johndoe",
			"expected_code":    http.StatusBadRequest,
			"expected_status":  "BAD REQUEST",
		},
		"USER_Register_VALIDATION_ERROR_email_invalid": {
			"request_name":     "John Doe",
			"request_email":    "johndoe",
			"request_password": "johndoe",
			"expected_code":    http.StatusBadRequest,
			"expected_status":  "BAD REQUEST",
		},
		"USER_Register_VALIDATION_ERROR_password_empty": {
			"request_name":     "John Doe",
			"request_email":    "johndoe@mail.com",
			"request_password": "",
			"expected_code":    http.StatusBadRequest,
			"expected_status":  "BAD REQUEST",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"name":"%s", "email":"%s", "password":"%s"}`,
				testItem["request_name"], testItem["request_email"], testItem["request_password"],
			)
			request := newRequest(http.MethodPost, registerUrl, requestBody)

			recoder := httptest.NewRecorder()
			app.ServeHTTP(recoder, request)
			response := recoder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			testResponse := new(TestResponse[any])
			json.Unmarshal(responseBody, testResponse)

			require.Equal(t, testItem["expected_code"].(int), testResponse.Code)
			require.Equal(t, testItem["expected_status"].(string), testResponse.Status)
		})

	}
}

func TestLogin(t *testing.T) {
	testItems := map[string]TestSchema{
		"USER_Login_OK": {
			"request_email":    "johndoe@mail.com",
			"request_password": "johndoe",
			"expected_code":    http.StatusOK,
			"expected_status":  "OK",
		},
		"USER_Login_VALIDATION_ERROR_email_empty": {
			"request_email":    "",
			"request_password": "johndoe",
			"expected_code":    http.StatusBadRequest,
			"expected_status":  "BAD REQUEST",
		},
		"USER_Login_VALIDATION_ERROR_email_doesnt_exists": {
			"request_email":    "abcdefg@mail.com",
			"request_password": "abcdefg",
			"expected_code":    http.StatusUnauthorized,
			"expected_status":  "UNAUTHORIZED",
		},
		"USER_Login_VALIDATION_ERROR_password_empty": {
			"request_email":    "johndoe@mail.com",
			"request_password": "",
			"expected_code":    http.StatusBadRequest,
			"expected_status":  "BAD REQUEST",
		},
		"USER_Login_VALIDATION_ERROR_password_wrong": {
			"request_email":    "johndoe@mail.com",
			"request_password": "abcdefg",
			"expected_code":    http.StatusUnauthorized,
			"expected_status":  "UNAUTHORIZED",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"email":"%s", "password":"%s"}`,
				testItem["request_email"], testItem["request_password"])
			request := newRequest(http.MethodPost, loginUrl, requestBody)

			recoder := httptest.NewRecorder()
			app.ServeHTTP(recoder, request)
			response := recoder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			testResponse := new(TestResponse[model.TokenData])
			json.Unmarshal(responseBody, testResponse)

			assert.Equal(t, testItem["expected_code"], testResponse.Code)
			assert.Equal(t, testItem["expected_status"], testResponse.Status)

			if testResponse.Code == http.StatusOK {
				validToken = testResponse.Data.AccessToken
				require.NotEmpty(t, validToken)

				var err error
				authData, err = utils.ParseAccessToken(viperConfig, validToken)
				require.Nil(t, err)

				accessTokenRedisKey := utils.GenerateAccessTokenRedisKey(authData.UserID)
				accessTokenRedis, err := redisClient.Get(context.Background(), accessTokenRedisKey).Result()
				require.Equal(t, accessTokenRedis, validToken)
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	post := model.PostResponse{
		Title:   "TEST_POST",
		Content: "TEST_CONTENT",
		Author:  authData.UserName,
	}

	testItems := map[string]TestSchema{
		"POST_Create_OK": {
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         validToken,
			"expected_code":         http.StatusOK,
			"expected_status":       "OK",
			"expected_data_title":   post.Title,
			"expected_data_content": post.Content,
			"expected_data_author":  post.Author,
		},
		"POST_Create_VALIDATION_ERROR_title_empty": {
			"request_title":         "",
			"request_content":       post.Content,
			"request_token":         validToken,
			"expected_code":         http.StatusBadRequest,
			"expected_status":       "BAD REQUEST",
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Create_VALIDATION_ERROR_content_empty": {
			"request_title":         post.Title,
			"request_content":       "",
			"request_token":         validToken,
			"expected_code":         http.StatusBadRequest,
			"expected_status":       "BAD REQUEST",
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Create_UNAUTHORIZED_ERROR_empty_token": {
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         "",
			"expected_code":         http.StatusUnauthorized,
			"expected_status":       "UNAUTHORIZED",
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Create_UNAUTHORIZED_ERROR_invalid_token": {
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         "this.is.invalid.token",
			"expected_code":         http.StatusUnauthorized,
			"expected_status":       "UNAUTHORIZED",
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"title":"%s", "content":"%s"}`,
				testItem["request_title"], testItem["request_content"])
			request := newRequestWithToken(http.MethodPost, postAdminUrl, requestBody, testItem["request_token"].(string))

			recorder := httptest.NewRecorder()
			app.ServeHTTP(recorder, request)
			response := recorder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			testResponse := new(TestResponse[model.PostResponse])

			require.Nil(t, json.Unmarshal(responseBody, testResponse))

			require.Equal(t, testItem["expected_code"].(int), testResponse.Code)
			require.Equal(t, testItem["expected_status"].(string), testResponse.Status)
			require.Equal(t, testItem["expected_data_title"].(string), testResponse.Data.Title)
			require.Equal(t, testItem["expected_data_content"].(string), testResponse.Data.Content)
			require.Equal(t, testItem["expected_data_author"].(string), testResponse.Data.Author)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	post := model.PostResponse{
		Title:   "UPDATE_TEST_POST",
		Content: "UPDATE_TEST_CONTENT",
		Author:  authData.UserName,
	}

	testItems := map[string]TestSchema{
		"POST_Update_OK": {
			"param_id":              "4",
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         validToken,
			"expected_code":         http.StatusOK,
			"expected_status":       "OK",
			"expected_data_id":      4,
			"expected_data_title":   post.Title,
			"expected_data_content": post.Content,
			"expected_data_author":  post.Author,
		},
		"POST_Update_VALIDATION_ERROR_title_empty": {
			"param_id":              "4",
			"request_title":         "",
			"request_content":       post.Content,
			"request_token":         validToken,
			"expected_code":         http.StatusBadRequest,
			"expected_status":       "BAD REQUEST",
			"expected_data_id":      0,
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Update_VALIDATION_ERROR_content_empty": {
			"param_id":              "4",
			"request_title":         post.Title,
			"request_content":       "",
			"request_token":         validToken,
			"expected_code":         http.StatusBadRequest,
			"expected_status":       "BAD REQUEST",
			"expected_data_id":      0,
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Update_UNAUTHORIZED_ERROR_empty_token": {
			"param_id":              "4",
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         "",
			"expected_code":         http.StatusUnauthorized,
			"expected_status":       "UNAUTHORIZED",
			"expected_data_id":      0,
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Update_UNAUTHORIZED_ERROR_invalid_token": {
			"param_id":              "4",
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         "this.is.invalid.token",
			"expected_code":         http.StatusUnauthorized,
			"expected_status":       "UNAUTHORIZED",
			"expected_data_id":      0,
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Update_NOT_FOUND_other_user_post": {
			"param_id":              "1",
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         validToken,
			"expected_code":         http.StatusNotFound,
			"expected_status":       "NOT FOUND",
			"expected_data_id":      0,
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
		"POST_Update_NOT_FOUND_not_existing_Post": {
			"param_id":              "1000",
			"request_title":         post.Title,
			"request_content":       post.Content,
			"request_token":         validToken,
			"expected_code":         http.StatusNotFound,
			"expected_status":       "NOT FOUND",
			"expected_data_id":      0,
			"expected_data_title":   "",
			"expected_data_content": "",
			"expected_data_author":  "",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestUrl := postAdminUrl + "/" + testItem["param_id"].(string)
			requestBody := fmt.Sprintf(`{"title":"%s", "content":"%s"}`,
				testItem["request_title"], testItem["request_content"])
			request := newRequestWithToken(http.MethodPut, requestUrl,
				requestBody, testItem["request_token"].(string))

			recorder := httptest.NewRecorder()
			app.ServeHTTP(recorder, request)
			response := recorder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			testResponse := new(TestResponse[model.PostResponse])

			require.Nil(t, json.Unmarshal(responseBody, testResponse))

			require.Equal(t, testItem["expected_code"].(int), testResponse.Code)
			require.Equal(t, testItem["expected_status"].(string), testResponse.Status)
			require.Equal(t, uint64(testItem["expected_data_id"].(int)), testResponse.Data.ID)
			require.Equal(t, testItem["expected_data_title"].(string), testResponse.Data.Title)
			require.Equal(t, testItem["expected_data_content"].(string), testResponse.Data.Content)
			require.Equal(t, testItem["expected_data_author"].(string), testResponse.Data.Author)
		})
	}
}

func TestDeletePost(t *testing.T) {
	testItems := map[string]TestSchema{
		"POST_Delete_UNAUTHORIZED_ERROR_empty_token": {
			"param_id":        "4",
			"request_token":   "",
			"expected_code":   http.StatusUnauthorized,
			"expected_status": "UNAUTHORIZED",
		},
		"POST_Delete_UNAUTHORIZED_ERROR_invalid_token": {
			"param_id":        "4",
			"request_token":   "this.is.invalid.token",
			"expected_code":   http.StatusUnauthorized,
			"expected_status": "UNAUTHORIZED",
		},
		"POST_Delete_NOT_FOUND_other_user_post": {
			"param_id":        "1",
			"request_token":   validToken,
			"expected_code":   http.StatusNotFound,
			"expected_status": "NOT FOUND",
		},
		"POST_Delete_NOT_FOUND_not_existing_Post": {
			"param_id":        "10000",
			"request_token":   validToken,
			"expected_code":   http.StatusNotFound,
			"expected_status": "NOT FOUND",
		},
		"POST_Delete_OK": {
			"param_id":        "4",
			"request_token":   validToken,
			"expected_code":   http.StatusOK,
			"expected_status": "OK",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestUrl := postAdminUrl + "/" + testItem["param_id"].(string)
			request := newRequestWithToken(http.MethodDelete, requestUrl, "", testItem["request_token"].(string))

			recorder := httptest.NewRecorder()
			app.ServeHTTP(recorder, request)
			response := recorder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			testResponse := new(TestResponse[any])

			require.Nil(t, json.Unmarshal(responseBody, testResponse))

			require.Equal(t, testItem["expected_code"].(int), testResponse.Code)
			require.Equal(t, testItem["expected_status"].(string), testResponse.Status)
		})
	}
}
