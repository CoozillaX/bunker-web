package rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KickPlayerRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 租赁服实体ID
	UID      int    `json:"uid" binding:"min=1"`       // 目标玩家的启动器UID
}

// KickPlayer godoc
//
//	@Summary		踢出租赁服玩家
//	@Description	踢出指定租赁服的指定玩家, 游戏账号需要拥有指定的租赁服, 且需要提供 API key
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			租赁服 (管理类)
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		KickPlayerRequest		true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/rental_server/kick_player [post]
func (*RentalServer) KickPlayer(c *gin.Context) {
	// Parse request
	var req KickPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	// Check helper if exist
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.MpayToken == "" {
		c.Error(giner.NewPublicGinError("请先绑定游戏账号"))
		return
	}
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser.MpayUser, nil)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Request player kick
	ginerr = g79.KickRentalServerPlayer(gu, req.ServerID, strconv.Itoa(req.UID))
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Return
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	// Create log
	c.Set("log", "租赁服玩家踢出成功")
}
