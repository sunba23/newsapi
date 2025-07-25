package database

import (
	"time"
)

type User struct {
	ID        string    `db:"id"`
	GoogleID  string    `db:"google_id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

type Tag struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type News struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Author    string    `db:"author"`
	CreatedAt time.Time `db:"created_at"`
	Tags      []Tag     `db:"-"`
}

type NewsWithTags struct {
	News
	TagID   *int    `db:"tag_id"`
	TagName *string `db:"tag_name"`
}
