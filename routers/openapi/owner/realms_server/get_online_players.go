package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetOnlinePlayersRequest struct {
	ServerID  string `json:"server_id" binding:"min=1"`                            // 山头服 SID
	Offset    int    `json:"offset" binding:"gte=0"`                               // 分页偏移
	Length    int    `json:"length" binding:"omitempty,gte=1,lte=100" example:"9"` // 分页长度, 默认为 9
	OrderType int    `json:"order_type" example:"0"`                               // 排序类型 (0, 1)
}

type GetOnlinePlayersResponse struct {
	giner.BasicResponse
	Data []map[string]any `json:"data"` // 在线玩家列表
}

// GetOnlinePlayers godoc
//
//	@Summary		获取山头服在线玩家列表
//	@Description	通过山头服 SID 分页查询当前在线玩家
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		GetOnlinePlayersRequest	true	"在线玩家查询参数"
//	@Success		200				{object}	GetOnlinePlayersResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/get_online_players [post]
func (*RealmsServer) GetOnlinePlayers(c *gin.Context) {
	var req GetOnlinePlayersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	if req.Length == 0 {
		req.Length = 9
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	players, ginerr := g79.QueryRealmsServerOnlinePlayers(gu, req.ServerID, req.Offset, req.Length, req.OrderType)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(players))
	c.Set("log", "查询山头服在线玩家成功")
}
