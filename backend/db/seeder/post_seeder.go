package seeder

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"backend/internal/utils"
	"time"

	"gorm.io/gorm"
)

func SeedsPost(db *gorm.DB, postRepository *repository.PostRepository, users *[]entity.User) ([]entity.Post, error) {
	var posts []entity.Post

	for _, user := range *users {
		tx := db.Begin()
		postCreated := &entity.Post{
			Title:     "Title" + utils.GenerateRandomString(10),
			Content:   "Content " + utils.GenerateRandomString(100),
			CreatedAt: time.Now(),
			UserID:    user.ID,
		}
		err := postRepository.Repository.Save(tx, postCreated)
		if err != nil {
			return posts, err
		}
		tx.Commit()

		posts = append(posts, *postCreated)
	}

	return posts, nil
}
