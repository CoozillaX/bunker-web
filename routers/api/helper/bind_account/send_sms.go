package bind_account

import (
	"bunker-web/models"
	"bunker-web/pkg/captcha"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/utils"
	"bunker-web/services/user"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SendSMSRequest struct {
	Mobile       string `json:"mobile" binding:"len=11"`
	CaptchaToken string `json:"captcha_token" binding:"min=1"`
}

type SendSMSResponseData struct {
	VerifyUrl string `json:"verify_url,omitempty"`
}

func (*BindAccount) SendSMS(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check captcha
	if !captcha.CheckTurnstileCaptchaToken(c.ClientIP(), req.CaptchaToken) {
		c.Error(giner.NewPublicGinError("验证码未通过"))
		return
	}
	// Check helper
	if usr.HelperMpayUser != nil {
		if usr.HelperMpayUser.GetToken() != "" {
			c.Error(giner.NewPublicGinError("创建失败, 已存在辅助用户账号"))
			return
		}
		if usr.HelperMpayUser.GetType() != models.MpayUserTypeAndroid {
			user.DeleteHelper(usr)
		}
	}
	// Check if has enough chances
	defer models.DBSave(usr)
	if !utils.IsToday(usr.LastGetSMSCodeAt.Time) {
		usr.SMSCodeTimes = 0
	}
	if usr.SMSCodeTimes > user.SMSCodeLimit {
		c.Error(giner.NewPublicGinError("今日获取验证码次数已达上限"))
		return
	}
	// Create helper user if not exist
	if usr.HelperMpayUser == nil {
		usr.HelperMpayUser = &models.AndroidMpayUser{}
	}
	// Try to request code
	if protocolErr := usr.HelperMpayUser.SMSLoginRequestCode(req.Mobile); protocolErr != nil {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetData(
			&SendSMSResponseData{
				VerifyUrl: protocolErr.VerifyUrl,
			},
		))
		return
	}
	// Renew usr sms info
	usr.SMSCodeTimes++
	usr.LastGetSMSCodeAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	// Return success
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage(
		fmt.Sprintf("获取成功, 您今天还可获取%d次手机验证码", user.SMSCodeLimit-usr.SMSCodeTimes),
	))
	// Create log
	c.Set("log", "获取helper手机登录验证码成功")
}
