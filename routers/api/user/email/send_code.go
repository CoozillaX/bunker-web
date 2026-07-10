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
	"strings"

	"github.com/gin-gonic/gin"
)

type SendCodeRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	ActionType   int    `json:"action_type" binding:"gte=0,lte=4"`
	CaptchaToken string `json:"captcha_token" binding:"min=1"`
}

// loadOptionalUser resolves the logged-in user when present.
// This route sits outside BearerHandler so reset-password can stay public;
// bind/unbind/change/remove still need a valid session from header or cookie.
func loadOptionalUser(c *gin.Context) *models.User {
	bearer := ""
	if v, ok := c.Get("bearer"); ok {
		if s, ok := v.(string); ok {
			bearer = s
		}
	}
	if bearer == "" {
		bearer = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	}
	if bearer == "" {
		if cookie, err := c.Cookie(sessions.SESSION_COOKIE_NAME); err == nil {
			bearer = cookie
		}
	}
	if bearer == "" {
		return nil
	}
	session, ok := sessions.GetSessionByBearer(bearer)
	if !ok {
		return nil
	}
	u, ok := session.Load("usr")
	if !ok {
		return nil
	}
	usr, ok := u.(*models.User)
	if !ok {
		return nil
	}
	return usr
}

func (*Email) SendCode(c *gin.Context) {
	usr := loadOptionalUser(c)
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
		if usr == nil {
			c.Error(giner.NewPublicGinError("无效请求"))
			return
		}
		if usr.Email != "" {
			c.Error(giner.NewPublicGinError("请先解绑当前邮箱"))
			return
		}
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
