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
	if usr.HelperMpayUser != nil && usr.HelperMpayUser.MpayToken != "" {
		c.Error(giner.NewPublicGinError("创建失败, 已存在辅助用户账号"))
		return
	}
	// Store to DB
	defer models.DBSave(usr.HelperMpayUser)
	// Try to login
	helper, ginerr := user.GetLoginHelperForHelper(usr)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if protocolErr := helper.GuestLogin(); protocolErr != nil {
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
