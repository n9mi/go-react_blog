package repository

import (
	"backend/internal/entity"
	"backend/internal/model"
	"backend/internal/utils"
	"strings"

	"gorm.io/gorm"
)

type PostRepository struct {
	Repository[entity.Post]
}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}

func (r *PostRepository) List(tx *gorm.DB, request *model.PostListRequest) ([]model.PostResponse, error) {
	var postList []model.PostResponse

	if request.Page > 0 && request.PageSize > 0 {
		tx = tx.Scopes(utils.Paginate(request.Page, request.PageSize))
	}

	query := tx.Model(new(entity.Post)).
		Order("posts.id asc").
		Select(`posts.id,
			posts.title,
			posts.content,
			posts.created_at,
			users.name as author`).
		Joins("inner join users on users.id = posts.user_id")

	if len(request.UserID) > 0 {
		query = query.Where("users.id = ?", request.UserID)
	}

	if len(request.TitleQuery) > 0 {
		query = query.Where("LOWER(posts.title) LIKE ?", "%"+strings.ToLower(request.TitleQuery)+"%")
	}

	if err := query.Scan(&postList).Error; err != nil {
		return postList, err
	}

	for i := range postList {
		postList[i].Content = postList[i].GetContentSummary()
	}

	return postList, nil
}

func (r *PostRepository) GetWithAuthor(tx *gorm.DB, post *model.PostResponse, ID uint64) error {
	return tx.Model(new(entity.Post)).
		Where("posts.id = ?", ID).
		Select(`posts.id,
			posts.title,
			posts.content,
			posts.created_at,
			users.name as author`).
		Joins("inner join users on users.id = posts.user_id").
		Scan(&post).Error
}

func (r *PostRepository) GetByIDandAuthorID(tx *gorm.DB, post *entity.Post, ID uint64, userID string) error {
	return tx.Where("id = ? and user_id = ?", ID, userID).First(post).Error
}
