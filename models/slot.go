package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Slot struct {
	gorm.Model
	UserID   uint
	GameID   int
	ExpireAt sql.NullTime
	Note     string
}
