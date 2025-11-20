package email

import (
	"bunker-web/configs"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/utils"
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gopkg.in/gomail.v2"
)

type verifyCode struct {
	Email     string
	Code      string
	SendTime  time.Time
	CheckTime int
}

var (
	//go:embed verify_template.html
	verify_template []byte
	verifyCodeCache *cache.Cache
)

func init() {
	verifyCodeCache = cache.New(time.Minute*10, time.Minute)
}

func SendHTMLEmail(to, subject, body string) error {
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("BunkerWeb <%s>", configs.GMAIL_ACCOUNT))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	// Set up the SMTP dialer
	d := gomail.NewDialer("smtp.gmail.com", 465, "", "")
	// Use OAuth2 authentication
	d.Auth = Auth()
	// Send the email
	return d.DialAndSend(m)
}

func SendVerifyEmail(username, action, to string) *gin.Error {
	var code *verifyCode
	// Check if the email is in cooldown
	if vc, ok := verifyCodeCache.Get(username + action); ok {
		code = vc.(*verifyCode)
		if time.Since(code.SendTime) < time.Second*50 {
			return giner.NewPublicGinError("发送验证码过于频繁")
		}
	} else {
		code = &verifyCode{
			Email: to,
			Code:  strings.ToUpper(utils.GenerateRandomString(6)),
		}
		verifyCodeCache.SetDefault(username+action, code)
	}
	code.SendTime = time.Now()
	// Format email body
	tmpl, err := template.New("verifyTemplate").Parse(string(verify_template))
	if err != nil {
		return giner.NewPrivateGinError(err)
	}
	data := map[string]string{
		"Username": username,
		"Code":     code.Code,
		"Action":   action,
		"Time":     code.SendTime.Format(time.DateTime),
	}
	output := bytes.NewBuffer([]byte{})
	if err = tmpl.Execute(output, data); err != nil {
		return giner.NewPrivateGinError(err)
	}
	return giner.NewPrivateGinError(SendHTMLEmail(to, "BunkerWeb 邮箱验证码", output.String()))
}

func CheckVerifyCode(username, action, email, code string) bool {
	// Toupper
	code = strings.ToUpper(code)
	// Check if the code is correct
	if vc, ok := verifyCodeCache.Get(username + action); ok {
		verifyCode := vc.(*verifyCode)
		// Check time
		verifyCode.CheckTime++
		if verifyCode.CheckTime > 5 {
			return false
		}
		// Check email
		if verifyCode.Email != email {
			return false
		}
		// Check code
		if verifyCode.Code != code {
			return false
		}
		// Check passed, Delete cache
		verifyCodeCache.Delete(username + action)
		return true
	}
	return false
}
