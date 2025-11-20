package email

import (
	"bunker-web/models"
	"bunker-web/pkg/email"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnbindRequest struct {
	EmailVerifyCode string `json:"email_verify_code" binding:"len=6"`
}

func (*Email) Unbind(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req UnbindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check email verify code
	if !email.CheckVerifyCode(usr.Username, email.EmailVerifyActionTypeMap[email.EmailVerifyActionTypeUnbind], usr.Email, req.EmailVerifyCode) {
		c.Error(giner.NewPublicGinError("无效的邮箱验证码"))
		return
	}
	// Update email
	usr.Email = ""
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("邮箱解绑成功"))
	// Create log
	c.Set("log", "邮箱解绑成功")
}
