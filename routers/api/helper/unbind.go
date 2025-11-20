package helper

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*Helper) UnBind(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.MpayToken == "" {
		c.Error(giner.NewPublicGinError("解绑失败, 不存在辅助用户"))
		return
	}
	// Delete helper
	if ginerr := user.DeleteHelper(usr); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("辅助用户解绑成功"))
	// Create log
	c.Set("log", "Helper解绑成功")
}
