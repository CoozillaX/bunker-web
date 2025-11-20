package email

import (
	"bunker-web/models"
	"bunker-web/pkg/captcha"
	"bunker-web/pkg/email"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/utils"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SendCodeRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	ActionType   int    `json:"action_type" binding:"gte=0,lte=4"`
	CaptchaToken string `json:"captcha_token" binding:"min=1"`
}

func (*Email) SendCode(c *gin.Context) {
	// Get data from session
	var usr *models.User
	{
		bearer, _ := c.Get("bearer")
		session, _ := sessions.GetSessionByBearer(bearer.(string))
		u, ok := session.Load("usr")
		if ok {
			usr = u.(*models.User)
		}
	}
	// Parse request
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check captcha token
	if !captcha.CheckTurnstileCaptchaToken(c.ClientIP(), req.CaptchaToken) {
		c.Error(giner.NewPublicGinError("验证码未通过"))
		return
	}
	// Send email by diffrent action type
	var to string
	switch req.ActionType {
	case email.EmailVerifyActionTypeBind:
		if !utils.IsValidEmail(req.Email) {
			c.Error(giner.NewPublicGinError("无效参数"))
			return
		}
		to = req.Email
	case email.EmailVerifyActionTypeUnbind, email.EmailVerifyActionTypeChangePassword, email.EmailVerifyActionTypeRemoveAccount:
		if usr == nil {
			c.Error(giner.NewPublicGinError("无效请求"))
			return
		}
		if usr.Email == "" {
			c.Error(giner.NewPublicGinError("请先绑定邮箱"))
			return
		}
		to = usr.Email
	case email.EmailVerifyActionTypeResetPassword:
		if req.Username == "" {
			c.Error(giner.NewPublicGinError("无效参数"))
			return
		}
		var ginerr *gin.Error
		usr, ginerr = user.QueryByUsername(req.Username)
		if ginerr != nil {
			c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("如果信息无误, 您将收到一封验证邮件"))
			return
		}
		if usr.Email == "" {
			c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("如果信息无误, 您将收到一封验证邮件"))
			return
		}
		to = usr.Email
	default:
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Send email
	if ginerr := email.SendVerifyEmail(usr.Username, email.EmailVerifyActionTypeMap[req.ActionType], to); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("如果信息无误, 您将收到一封验证邮件"))
	// Create log
	c.Set("log", "请求发送邮箱验证码")
}
