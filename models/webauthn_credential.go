package models

import (
	"gorm.io/gorm"
)

type WebAuthnCredential struct {
	gorm.Model
	RawID  string
	UserID uint
	Base64 string
}
