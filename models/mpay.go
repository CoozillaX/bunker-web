package models

import (
	"bunker-core/protocol/defines"

	"gorm.io/gorm"
)

type MpayUser struct {
	gorm.Model
	*defines.MpayUser
}
