package bind_account

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MobileRequest struct {
	Mobile  string `json:"mobile" binding:"len=11"`
	SMSCode string `json:"smscode" binding:"omitempty,len=6"`
}

type MobileResponseData struct {
	VerifyUrl string `json:"verify_url,omitempty"`
}

func (*BindAccount) Mobile(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req MobileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check mpay user status
	if usr.OwnerMpayUser == nil {
		c.Error(giner.NewPublicGinError("请先获取手机验证码"))
		return
	} else if usr.OwnerMpayUser.GetToken() != "" {
		c.Error(giner.NewPublicGinError("绑定失败, 已绑定游戏账号"))
		return
	}
	// Create helper user if not exist
	if usr.OwnerMpayUser == nil {
		usr.OwnerMpayUser = &models.AndroidMpayUser{}
	}
	defer models.DBSave(usr)
	// Try to login
	if protocolErr := usr.OwnerMpayUser.SMSLoginVerifyCode(req.Mobile, req.SMSCode); protocolErr != nil {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).
			SetMessage(protocolErr.Message).
			SetData(
				&MobileResponseData{
					VerifyUrl: protocolErr.VerifyUrl,
				},
			))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("游戏账号绑定成功"))
	// Create log
	c.Set("log", "绑定Owner成功(手机)")
}
