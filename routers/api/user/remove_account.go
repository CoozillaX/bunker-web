package user

import (
	"bunker-web/models"
	"bunker-web/pkg/email"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RemoveAccountRequest struct {
	EmailVerifyCode string `json:"email_verify_code" binding:"len=6"`
}

func (*User) RemoveAccount(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req RemoveAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check email verify code
	if !email.CheckVerifyCode(usr.Username, email.EmailVerifyActionTypeMap[email.EmailVerifyActionTypeRemoveAccount], usr.Email, req.EmailVerifyCode) {
		c.Error(giner.NewPublicGinError("无效的邮箱验证码"))
		return
	}
	// Remove account
	if ginerr := user.Remove(usr); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Remove session
	sessions.DeleteSessionByBearer(bearer.(string))
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("您已成功删除账户, 感谢使用!"))
	// Create log
	c.Set("log", "删除账户成功")
}
