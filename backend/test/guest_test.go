package test

import (
	"backend/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	postGuestUrl = "http://127.0.0.1:5000/api/posts"
	postAdminUrl = "http://127.0.0.1:5000/api/admin/posts"
)

func TestPostList(t *testing.T) {
	request := newRequest(http.MethodGet, postGuestUrl, "")

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)
	response := recorder.Result()

	responseBody, _ := io.ReadAll(response.Body)
	resultMap := make(map[string]interface{})

	// Test if response can be parsed as map
	require.Equal(t, nil, json.Unmarshal(responseBody, &resultMap))

	// Test if response.code is 200
	require.Equal(t, 200, int(resultMap["code"].(float64)))
	// Test if response.status is OK
	require.Equal(t, "OK", resultMap["status"].(string))
	// Test if response.data is contains item that equal to seeder count
	responseList := resultMap["data"].([]interface{})
	require.True(t, len(responseList) > 0)
}

func TestPostGetByID(t *testing.T) {
	testItems := map[string]TestSchema{
		"POST_GetByID_OK": {
			"param_id":         "1",
			"code":             http.StatusOK,
			"status":           "OK",
			"expected_data_id": 1,
		},
		"POST_GetByID_NOTFOUND_ID=1000": {
			"param_id":         "1000",
			"code":             http.StatusNotFound,
			"status":           "NOT FOUND",
			"expected_data_id": 0,
		},
		"POST_GetByID_NOTFOUND_ID=abc": {
			"param_id":         "abc",
			"code":             http.StatusNotFound,
			"status":           "NOT FOUND",
			"expected_data_id": 0,
		},
		"POST_GetByID_NOTFOUND_ID=1abc": {
			"param_id":         "1abc",
			"code":             http.StatusNotFound,
			"status":           "NOT FOUND",
			"expected_data_id": 0,
		},
	}

	for testName, testMap := range testItems {
		t.Run(testName, func(t *testing.T) {
			request := newRequest(http.MethodGet, postGuestUrl+"/"+testMap["param_id"].(string), "")

			recorder := httptest.NewRecorder()
			app.ServeHTTP(recorder, request)
			response := recorder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			testResponse := new(TestResponse[model.PostResponse])

			require.Equal(t, nil, json.Unmarshal(responseBody, &testResponse))
			require.Equal(t, testMap["code"].(int), testResponse.Code)
			require.Equal(t, testMap["status"].(string), testResponse.Status)
			require.Equal(t, uint64(testMap["expected_data_id"].(int)), testResponse.Data.ID)
		})
	}
}
