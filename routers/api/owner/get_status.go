package owner

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetStatusResponseData struct {
	Set      bool   `json:"set"`
	Username string `json:"username,omitempty"`
	Lv       int    `json:"lv,omitempty"`
	Exp      int    `json:"exp"`
	TotalExp int    `json:"total_exp,omitempty"`
}

func (*Owner) GetStatus(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check owner if exist
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.MpayToken == "" {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).
			SetMessage("未绑定游戏账号").
			SetData(&GetStatusResponseData{
				Set: false,
			}),
		)
		// Create log
		c.Set("log", "未绑定游戏账号")
		return
	}
	// We don't need to login again if x19 user exists in session
	defer models.DBSave(usr.OwnerMpayUser)
	// Relogin
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser.MpayUser, nil)
	if ginerr != nil {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).
			SetMessage(giner.GetPublicErrorString(ginerr)).
			SetData(&GetStatusResponseData{
				Set: true,
			}),
		)
		// Create log
		c.Set("log", giner.GetPublicErrorString(ginerr))
		return
	}
	// Get helper grow info
	lv, exp, totalExp, ginerr := g79.GetG79LauncherLevel(gu)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Pack x19 user info when login success
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(
		&GetStatusResponseData{
			Set:      true,
			Username: gu.Username,
			Lv:       lv,
			Exp:      exp,
			TotalExp: totalExp,
		},
	))
	// Create log
	c.Set("log", "获取Owner信息成功")
}
