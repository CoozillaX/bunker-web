package bind_account

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetQRCodeResponseData struct {
	UUID      string `json:"uuid,omitempty"`
	ImageData []byte `json:"image_data,omitempty"`
	VerifyUrl string `json:"verify_url,omitempty"`
}

func (*BindAccount) GetQRCode(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check helper
	if usr.OwnerMpayUser != nil {
		if usr.OwnerMpayUser.GetToken() != "" {
			c.Error(giner.NewPublicGinError("无法获取二维码, 已存在辅助用户账号"))
			return
		}
		if usr.OwnerMpayUser.GetType() != models.MpayUserTypeWindows {
			user.DeleteOwner(usr)
		}
	}
	// Create helper user if not exist
	if usr.OwnerMpayUser == nil {
		usr.OwnerMpayUser = &models.WindowsMpayUser{}
	}
	defer models.DBSave(usr)
	// Try to request code
	uuid, image, protocolErr := usr.OwnerMpayUser.QRCodeLoginGetUUID()
	if protocolErr != nil {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetData(
			&GetQRCodeResponseData{
				VerifyUrl: protocolErr.VerifyUrl,
			},
		))
		return
	}
	// Return success
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(
		&GetQRCodeResponseData{
			UUID:      uuid,
			ImageData: image,
		},
	))
	// Create log
	c.Set("log", "获取Owner登录二维码成功")
}
