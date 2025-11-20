package models

import (
	"time"
)

type Log struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	IPAddress    string
	Method       string
	Path         string
	UserID       uint
	PublicError  string
	PrivateError string
	Message      string
}
