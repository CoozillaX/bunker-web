package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlocklistRemoveRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
	UID      string `json:"uid" binding:"min=1"`       // 玩家 UID
}

// BlocklistRemove godoc
//
//	@Summary		删除山头服黑名单
//	@Description	将指定玩家从山头服黑名单移除; 玩家不在黑名单时操作仍会成功
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		BlocklistRemoveRequest	true	"山头服 SID 和玩家 UID"
//	@Success		200				{object}	giner.BasicResponse			"成功时返回"
//	@Router			/openapi/owner/realms_server/blocklist_remove [post]
func (*RealmsServer) BlocklistRemove(c *gin.Context) {
	var req BlocklistRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.RemoveRealmsServerBlocklistPlayer(c.Request.Context(), gu, req.ServerID, req.UID); ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "删除山头服黑名单成功")
}
