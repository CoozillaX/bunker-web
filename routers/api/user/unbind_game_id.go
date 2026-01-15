package user

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*User) UnbindGameId(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check if binded
	if usr.GameID == 0 {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetMessage("当前未绑定游戏ID, 无需解绑"))
		// Create log
		c.Set("log", "无法解绑未绑定的游戏ID")
		return
	}
	// Set game id
	usr.GameID = 0
	// Save user
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("解绑成功"))
	// Create log
	c.Set("log", "解绑成功")
}
