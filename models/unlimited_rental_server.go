package models

import "gorm.io/gorm"

type UnlimitedRentalServer struct {
	gorm.Model
	OperatorID uint
	ServerCode string
}
