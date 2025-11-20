package owner

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UseGiftCodeRequest struct {
	Code string `json:"code" binding:"len=10"`
}

func (*Owner) UseGiftCode(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req UseGiftCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check owner if exist
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.MpayToken == "" {
		c.Error(giner.NewPublicGinError("未创建游戏账号"))
		return
	}
	// Store to DB
	defer models.DBSave(usr.OwnerMpayUser)
	// Relogin
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser.MpayUser, nil)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Use gift code
	if ginerr := g79.UseGiftCode(gu, req.Code); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("兑换码使用成功"))
	// Create log
	c.Set("log", "成功使用兑换码")
}
