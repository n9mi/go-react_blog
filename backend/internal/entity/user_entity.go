package entity

import (
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	Posts     []Post
}

func (e *User) EntityName() string {
	return "user"
}
