package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type UserBanRecord struct {
	gorm.Model
	UserID uint
	Until  sql.NullTime
	Reason string
}
