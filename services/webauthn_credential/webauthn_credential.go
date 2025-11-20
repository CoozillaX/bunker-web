package webauthn_credential

import (
	"bunker-web/models"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/go-webauthn/webauthn/webauthn"
)

func convertToModel(credential *webauthn.Credential) *models.WebAuthnCredential {
	// Json marshal
	bytes, _ := json.Marshal(credential)
	// Base64 encode
	return &models.WebAuthnCredential{
		RawID:  strings.TrimRight(base64.RawURLEncoding.EncodeToString(credential.ID), "="),
		Base64: base64.StdEncoding.EncodeToString(bytes),
	}
}

func convertToRaw(credential *models.WebAuthnCredential) *webauthn.Credential {
	// Base64 decode
	bytes, _ := base64.StdEncoding.DecodeString(credential.Base64)
	// Json unmarshal
	var result webauthn.Credential
	json.Unmarshal(bytes, &result)
	return &result
}

func StoreToDB(credential *webauthn.Credential, userID uint) error {
	credentialm := convertToModel(credential)
	credentialm.UserID = userID
	query := models.DB.Create(credentialm)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func Remove(credential *models.WebAuthnCredential) error {
	return models.DBRemove(credential)
}

func QueryModelByID(id uint) (*models.WebAuthnCredential, error) {
	var result models.WebAuthnCredential
	query := models.DB.First(&result, id)
	if query.Error != nil {
		return nil, query.Error
	}
	return &result, nil
}

func QueryModelByRawID(rawid []byte) (*models.WebAuthnCredential, error) {
	var result models.WebAuthnCredential
	query := models.DB.Where("raw_id = ?", strings.TrimRight(base64.RawURLEncoding.EncodeToString(rawid), "=")).First(&result)
	if query.Error != nil {
		return nil, query.Error
	}
	return &result, nil
}

func QueryModelsByUserID(userID uint) ([]models.WebAuthnCredential, error) {
	var result []models.WebAuthnCredential
	query := models.DB.Where("user_id = ?", userID).Find(&result)
	if query.Error != nil {
		return nil, query.Error
	}
	return result, nil
}

func QueryRawsByUserID(userID uint) ([]webauthn.Credential, error) {
	var credentialms []models.WebAuthnCredential
	query := models.DB.Where("user_id = ?", userID).Find(&credentialms)
	if query.Error != nil {
		return nil, query.Error
	}
	var result []webauthn.Credential
	for _, credentialm := range credentialms {
		result = append(result, *convertToRaw(&credentialm))
	}
	return result, nil
}
