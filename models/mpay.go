package models

import (
	"bunker-core/protocol/mpay"
	"bunker-core/protocol/mpay/android"
	"bunker-core/protocol/mpay/windows"

	"gorm.io/gorm"
)

type MpayUserType string

const (
	MpayUserTypeAndroid MpayUserType = "android"
	MpayUserTypeWindows MpayUserType = "windows"
)

type MpayUser interface {
	mpay.MpayUser
	GetID() uint
	GetType() MpayUserType
}

type AndroidMpayUser struct {
	gorm.Model
	android.AndroidMpayUser
}

func (a *AndroidMpayUser) GetID() uint {
	return a.ID
}

func (a *AndroidMpayUser) GetType() MpayUserType {
	return MpayUserTypeAndroid
}

type WindowsMpayUser struct {
	gorm.Model
	windows.WindowsMpayUser
}

func (w *WindowsMpayUser) GetID() uint {
	return w.ID
}

func (w *WindowsMpayUser) GetType() MpayUserType {
	return MpayUserTypeWindows
}
