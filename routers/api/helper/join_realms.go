package helper

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JoinRealmsRequest struct {
	Code string `json:"code" binding:"min=1,max=20"`
}

func (*Helper) JoinRealms(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req JoinRealmsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("未创建辅助用户"))
		return
	}
	// Store to DB
	defer models.DBSave(usr.HelperMpayUser)
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Join realms server
	if ginerr := gu.JoinRealmsServer(req.Code); ginerr != nil {
		c.Error(giner.NewPublicGinError(ginerr.Error()))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("已成功加入山头服"))
	// Create log
	c.Set("log", "加入山头服成功")
}
