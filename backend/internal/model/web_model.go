package model

type DataResponse[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type MessagesResponse struct {
	Code     int      `json:"code"`
	Status   string   `json:"status"`
	Messages []string `json:"messages"`
}
