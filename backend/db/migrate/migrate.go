package migrate

import (
	"backend/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Post{}); err != nil {
		return err
	}

	return nil
}

func Drop(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&entity.User{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Post{}); err != nil {
		return err
	}

	return nil
}
