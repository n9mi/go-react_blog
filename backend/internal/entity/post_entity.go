package entity

import "time"

type Post struct {
	ID        uint64 `gorm:"primaryKey"`
	Title     string
	Content   string
	CreatedAt time.Time `gorm:"<-create"`
	UserID    string
}

func (e *Post) EntityName() string {
	return "post"
}
