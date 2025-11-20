package user

import (
	"bunker-web/pkg/captcha"
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	UserName     string `json:"username" binding:"min=1"`
	Password     string `json:"password" binding:"len=64"`
	CaptchaToken string `json:"captcha_token" binding:"min=1"`
}

func (*User) Register(c *gin.Context) {
	// Parse request
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check captcha
	if !captcha.CheckTurnstileCaptchaToken(c.ClientIP(), req.CaptchaToken) {
		c.Error(giner.NewPublicGinError("验证码未通过"))
		return
	}
	// Create user
	if _, ginerr := user.Create(req.UserName, req.Password, user.PermissionGuest); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("注册成功"))
	// Create log
	c.Set("log", "注册成功")
}
