package bind_account

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QRCodeRequest struct {
	UUID string `json:"uuid" binding:"omitempty,gte=1"`
}

type QRCodeResponseData struct {
	VerifyUrl string `json:"verify_url,omitempty"`
}

func (*BindAccount) QRCode(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check mpay user status
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.GetType() != models.MpayUserTypeWindows {
		c.Error(giner.NewPublicGinError("请先获取登录二维码"))
		return
	}
	if usr.OwnerMpayUser.GetToken() != "" {
		c.Error(giner.NewPublicGinError("绑定失败, 已存在辅助用户账号"))
		return
	}
	// Store to DB
	defer models.DBSave(usr)
	// Try to login
	if protocolErr := usr.OwnerMpayUser.QRCodeLoginByUUID(req.UUID); protocolErr != nil {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).
			SetMessage(protocolErr.Message).
			SetData(
				&QRCodeResponseData{
					VerifyUrl: protocolErr.VerifyUrl,
				},
			))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("辅助用户绑定成功"))
	// Create log
	c.Set("log", "绑定Owner成功(二维码)")
}
