package model

import "strings"

type PostListRequest struct {
	Page       int
	PageSize   int
	UserID     string
	TitleQuery string
}

type PostCreateRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	AuthorID string `json:"-" validate:"required"`
}

type PostUpdateRequest struct {
	ID       uint64 `json:"-" validate:"required,min=1"`
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	AuthorID string `json:"-" validate:"required"`
}

type PostResponse struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Author    string `json:"author"`
}

// GetContentSummary returns summary of the Content by returning first 50 words
// If content has less than 50 words, GetContentSummary returns the whole Content
func (e *PostResponse) GetContentSummary() string {
	splitted := strings.Split(e.Content, " ")

	wordToTrim := len(splitted)
	if wordToTrim < 50 {
		return e.Content
	}
	wordToTrim = 50

	return strings.Join(splitted[:wordToTrim], " ")
}

type PostGetByIDRequest struct {
	ID uint64
}

type PostDeleteRequest struct {
	ID     uint64 `validate:"required,min=1"`
	UserID string `validate:"required"`
}
