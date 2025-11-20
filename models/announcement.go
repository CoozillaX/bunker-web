package models

import (
	"gorm.io/gorm"
)

type Announcement struct {
	gorm.Model
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorID   *uint  `json:"-"`
	Author     *User  `json:"-"`
	AuthorName string `json:"author_name" gorm:"-"`
}
