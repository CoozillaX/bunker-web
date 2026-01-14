package slot

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/slot"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetGameIDRequest struct {
	ID         uint   `json:"id" binding:"gt=0"`
	ServerCode string `json:"server_code" binding:"min=1,max=20"`
}

func (*Slot) SetGameID(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req SetGameIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Query slot
	s, err := slot.QueryByID(req.ID)
	if err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check permission
	if s.UserID != usr.ID {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check if binded game id
	if usr.GameID == 0 {
		c.Error(giner.NewPublicGinError("请先绑定游戏ID"))
		return
	}
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		c.JSON(http.StatusOK, giner.MakeHTTPResponse(false).SetMessage("未创建辅助用户"))
		c.Set("log", "未创建Helper")
		return
	}
	// g79 login
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser)
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
	// Check if equal to game id
	if usr.GameID == partialServerInfo.OwnerID {
		c.Error(giner.NewPublicGinError("此服务器拥有者的ID与当前账户绑定的游戏ID一致, 无需重复绑定"))
		return
	}
	// Set game id
	if ginerr := slot.SetGameID(s, usr.ID, partialServerInfo.OwnerID, fmt.Sprintf("已绑定服务器 %s 拥有者的游戏ID", req.ServerCode)); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("设置成功"))
	// Create log
	c.Set("log", "设置 slot 游戏ID成功")
}
