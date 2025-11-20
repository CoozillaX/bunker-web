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
	Username         string
	Password         string
	Email            string
	Permission       uint
	UnlimitedUntil   sql.NullTime
	GameID           int
	ExpireAt         sql.NullTime
	HelperMpayUserID *uint
	HelperMpayUser   *MpayUser `gorm:"constraint:OnDelete:SET NULL,OnUpdate:RESTRICT;"`
	OwnerMpayUserID  *uint
	OwnerMpayUser    *MpayUser `gorm:"constraint:OnDelete:SET NULL,OnUpdate:RESTRICT;"`
	SMSCodeTimes     int
	LastGetSMSCodeAt sql.NullTime
	APIKey           string
	ResponseTo       string // PhoenixBuilder specific field
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
