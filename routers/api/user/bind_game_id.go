package user

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (*User) BindGameId(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check if binded
	if usr.GameID != 0 {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetMessage("当前已绑定游戏ID, 请先解绑后再绑定"))
		// Create log
		c.Set("log", "无法重复绑定游戏ID")
		return
	}
	// Check owner if exist
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.GetToken() == "" {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetMessage("请登录游戏账号后再进行游戏ID绑定"))
		// Create log
		c.Set("log", "未登录游戏账号, 无法绑定游戏ID")
		return
	}
	// g79 login
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Set game id
	usr.GameID, _ = strconv.Atoi(gu.EntityID)
	// Save user
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("绑定成功"))
	// Create log
	c.Set("log", "绑定成功")
}
