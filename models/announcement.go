package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Announcement struct {
	gorm.Model
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	PinnedAt   sql.NullTime `json:"-"`
	AuthorID   *uint        `json:"-"`
	Author     *User        `json:"-"`
	AuthorName string       `json:"author_name" gorm:"-"`
}
