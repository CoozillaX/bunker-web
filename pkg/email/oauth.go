package email

import (
	"bunker-web/configs"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func getAccessToken() (string, error) {
	tokenURL := "https://accounts.google.com/o/oauth2/token"
	// request body
	data := url.Values{}
	data.Set("client_id", configs.GMAIL_CLIENT_ID)
	data.Set("client_secret", configs.GMAIL_CLIENT_SECRET)
	data.Set("refresh_token", configs.GMAIL_REFRESH_TOKEN)
	data.Set("grant_type", "refresh_token")
	// send request
	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	// check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get access token, status code: %d", resp.StatusCode)
	}
	// decode response
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}
	return tokenResponse.AccessToken, nil
}

type SMTPOAuth2 struct {
	Username string
}

func Auth() smtp.Auth {
	return &SMTPOAuth2{
		Username: configs.GMAIL_ACCOUNT,
	}
}

func (a *SMTPOAuth2) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "XOAUTH2", []byte(""), nil
}

func (a *SMTPOAuth2) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// Get access token
		accessToken, err := getAccessToken()
		if err != nil {
			return nil, fmt.Errorf("error getting access token: %v", err)
		}
		// Format token
		authToken := fmt.Sprintf("user=%s\001auth=Bearer %s\001\001", a.Username, accessToken)
		return []byte(authToken), nil
	}
	return nil, nil
}
