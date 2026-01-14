package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username           string
	Password           string
	Email              string
	Permission         uint
	UnlimitedUntil     sql.NullTime
	GameID             int
	ExpireAt           sql.NullTime
	HelperMpayUserID   uint
	HelperMpayUserType MpayUserType
	OwnerMpayUserID    uint
	OwnerMpayUserType  MpayUserType
	SMSCodeTimes       int
	LastGetSMSCodeAt   sql.NullTime
	APIKey             string
	HelperMpayUser     MpayUser `gorm:"-"`
	OwnerMpayUser      MpayUser `gorm:"-"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.HelperMpayUser != nil {
		if err := tx.Save(u.HelperMpayUser).Error; err != nil {
			return err
		}
		u.HelperMpayUserID = u.HelperMpayUser.GetID()
		u.HelperMpayUserType = u.HelperMpayUser.GetType()
	} else {
		u.HelperMpayUserID = 0
		u.HelperMpayUserType = ""
	}

	if u.OwnerMpayUser != nil {
		if err := tx.Save(u.OwnerMpayUser).Error; err != nil {
			return err
		}
		u.OwnerMpayUserID = u.OwnerMpayUser.GetID()
		u.OwnerMpayUserType = u.OwnerMpayUser.GetType()
	} else {
		u.OwnerMpayUserID = 0
		u.OwnerMpayUserType = ""
	}
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if u.HelperMpayUserID != 0 && u.HelperMpayUserType != "" {
		switch u.HelperMpayUserType {
		case MpayUserTypeAndroid:
			var am AndroidMpayUser
			if err := tx.First(&am, u.HelperMpayUserID).Error; err == nil {
				u.HelperMpayUser = &am
			}
		case MpayUserTypeWindows:
			var wm WindowsMpayUser
			if err := tx.First(&wm, u.HelperMpayUserID).Error; err == nil {
				u.HelperMpayUser = &wm
			}
		}
	}

	if u.OwnerMpayUserID != 0 && u.OwnerMpayUserType != "" {
		switch u.OwnerMpayUserType {
		case MpayUserTypeAndroid:
			var am AndroidMpayUser
			if err := tx.First(&am, u.OwnerMpayUserID).Error; err == nil {
				u.OwnerMpayUser = &am
			}
		case MpayUserTypeWindows:
			var wm WindowsMpayUser
			if err := tx.First(&wm, u.OwnerMpayUserID).Error; err == nil {
				u.OwnerMpayUser = &wm
			}
		}
	}
	return nil
}

func convertToRaw(credential *WebAuthnCredential) *webauthn.Credential {
	// Base64 decode
	bytes, _ := base64.StdEncoding.DecodeString(credential.Base64)
	// Json unmarshal
	var result webauthn.Credential
	json.Unmarshal(bytes, &result)
	return &result
}

func (u *User) WebAuthnID() []byte {
	return binary.BigEndian.AppendUint32([]byte{}, uint32(u.ID))
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Username
}

func (u *User) WebAuthnCredentials() (retval []webauthn.Credential) {
	var credentialms []WebAuthnCredential
	if query := DB.Where("user_id = ?", u.ID).Find(&credentialms); query.Error != nil {
		return nil
	}
	for _, credentialm := range credentialms {
		retval = append(retval, *convertToRaw(&credentialm))
	}
	return
}

func (u *User) WebAuthnIcon() string {
	return ""
}
