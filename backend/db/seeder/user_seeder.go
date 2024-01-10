package seeder

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Seeds `users` table by populating dummy user such as user1@mail.com, etc.
func SeedsUser(db *gorm.DB, userRepository *repository.UserRepository) ([]entity.User, error) {
	var users []entity.User

	for i := 1; i <= 3; i++ {
		userPassword, _ := utils.HashUserPassword(fmt.Sprintf("user%d", i))

		tx := db.Begin()
		userCreated := entity.User{
			ID:        "USR-" + utils.GenerateRandomString(20),
			Name:      fmt.Sprintf("user %d", i),
			Email:     fmt.Sprintf("user%d@mail.com", i),
			Password:  userPassword,
			CreatedAt: time.Now(),
		}
		err := userRepository.Repository.Save(tx, &userCreated)
		if err != nil {
			return users, err
		}
		tx.Commit()

		users = append(users, userCreated)
	}

	return users, nil
}
