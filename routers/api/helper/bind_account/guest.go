package bind_account

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GuestResponseData struct {
	VerifyUrl string `json:"verify_url,omitempty"`
}

func (*BindAccount) Guest(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check helper
	if usr.HelperMpayUser != nil {
		if usr.HelperMpayUser.GetToken() != "" {
			c.Error(giner.NewPublicGinError("无法获取二维码, 已存在辅助用户账号"))
			return
		}
		if usr.HelperMpayUser.GetType() != models.MpayUserTypeAndroid {
			user.DeleteHelper(usr)
		}
	}
	// Create helper user if not exist
	if usr.HelperMpayUser == nil {
		usr.HelperMpayUser = &models.AndroidMpayUser{}
	}
	defer models.DBSave(usr)
	// Try to login
	if protocolErr := usr.HelperMpayUser.GuestLogin(); protocolErr != nil {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).
			SetMessage(protocolErr.Message).
			SetData(
				&GuestResponseData{
					VerifyUrl: protocolErr.VerifyUrl,
				},
			))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("辅助用户绑定成功"))
	// Create log
	c.Set("log", "新建Helper成功")
}
