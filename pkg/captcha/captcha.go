package captcha

import (
	"bunker-web/configs"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func CheckTurnstileCaptchaToken(ip, token string) bool {
	reqForm := url.Values{
		"secret":   {configs.CAPTCHA_SECRET_KEY},
		"response": {token},
		"remoteip": {ip},
	}
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", reqForm)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var data struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return false
	}
	return data.Success
}
