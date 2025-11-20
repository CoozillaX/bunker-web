package user

import (
	"bunker-web/configs"
	"bunker-web/models"
	"bunker-web/pkg/email"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/utils"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	Username        string `json:"username" binding:"min=1"`
	EmailVerifyCode string `json:"email_verify_code" binding:"len=6"`
	NewPassword     string `json:"new_password" binding:"len=64"`
}

func (*User) ResetPassword(c *gin.Context) {
	// Parse request
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	if req.NewPassword == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user by username
	usr, ginerr := user.QueryByUsername(req.Username)
	if ginerr != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check email verify code
	if !email.CheckVerifyCode(usr.Username, email.EmailVerifyActionTypeMap[email.EmailVerifyActionTypeResetPassword], usr.Email, req.EmailVerifyCode) {
		c.Error(giner.NewPublicGinError("无效的邮箱验证码"))
		return
	}
	// Update password
	usr.Password = utils.SHA256Hex([]byte(req.NewPassword + configs.USER_PSW_SALT))
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("密码重置成功"))
	// Create log
	c.Set("log", "密码重置成功")
}
