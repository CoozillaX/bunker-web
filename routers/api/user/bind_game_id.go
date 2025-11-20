package user

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BindGameIdRequest struct {
	ServerCode string `json:"server_code" binding:"min=1,max=20"`
}

func (*User) BindGameId(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req BindGameIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check if binded
	if usr.GameID != 0 {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetMessage("无法重复绑定游戏ID"))
		// Create log
		c.Set("log", "无法重复绑定游戏ID")
		return
	}
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.MpayToken == "" {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetMessage("未创建辅助用户"))
		// Create log
		c.Set("log", "未创建Helper")
		return
	}
	// g79 login
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser.MpayUser, nil)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Query server base info
	partialServerInfo, ginerr := g79.QueryRentalServer(gu, req.ServerCode)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Set game id
	usr.GameID = partialServerInfo.OwnerID
	// Save user
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("绑定成功"))
	// Create log
	c.Set("log", "绑定成功")
}
