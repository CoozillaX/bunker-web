package models

import (
	"gorm.io/gorm"
)

type RedeemCode struct {
	gorm.Model
	Code     string `gorm:"unique"`
	CodeType int
	Used     bool
	UserID   uint
	Note     string
}
