package rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BanPlayerRequest struct {
	ServerID string `json:"server_id" binding:"min=1"`  // 租赁服实体ID
	UID      int    `json:"uid" binding:"min=1"`        // 目标玩家的启动器UID
	Status   int    `json:"status" binding:"oneof=0 1"` // 是否为黑名单状态 (0, 1)
}

// BanPlayer godoc
//
//	@Summary		更新租赁服玩家黑名单状态
//	@Description	更新租赁服玩家黑名单状态, 游戏账号需要拥有指定的租赁服, 且需要提供 API key
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			租赁服 (管理类)
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		BanPlayerRequest		true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/rental_server/ban_player [post]
func (*RentalServer) BanPlayer(c *gin.Context) {
	// Parse request
	var req BanPlayerRequest
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
	// Request player ban
	ginerr = g79.UpdateRentalServerPlayerBanStatus(gu, req.ServerID, strconv.Itoa(req.UID), req.Status)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Return
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	// Create log
	c.Set("log", "目标玩家黑名单状态更新成功")
}
