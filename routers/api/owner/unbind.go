package owner

import (
	"bunker-web/models"
	"bunker-web/services/user"
	"net/http"

	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"

	"github.com/gin-gonic/gin"
)

func (*Owner) UnBind(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check owner if exist
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.MpayToken == "" {
		c.Error(giner.NewPublicGinError("解绑失败, 未绑定游戏账号"))
		return
	}
	// Delete owner
	if ginerr := user.DeleteOwner(usr); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("游戏账号解绑成功"))
	// Create log
	c.Set("log", "Owner解绑成功")
}
