package fbtoken

import (
	"bunker-web/configs"
	"bunker-web/pkg/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encrypt(username, password string) (fbtoken string, err error) {
	jsonBytes, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	fbtokenBytes, err := utils.AES_256_CFBEncrypt(configs.FBTOKEN_KEY, jsonBytes, configs.FBTOKEN_IV)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fbtokenBytes), nil
}

func Decrypt(fbtoken string) (username, password string, err error) {
	fbtokenBytes, err := base64.StdEncoding.DecodeString(fbtoken)
	if err != nil {
		return "", "", err
	}
	jsonBytes, err := utils.AES_256_CFBDecrypt(configs.FBTOKEN_KEY, fbtokenBytes, configs.FBTOKEN_IV)
	if err != nil {
		return "", "", err
	}
	var result struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return "", "", err
	}
	if result.Username == "" || result.Password == "" {
		return "", "", fmt.Errorf("invaild token")
	}
	return result.Username, result.Password, nil
}
