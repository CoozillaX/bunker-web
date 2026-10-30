package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlocklistQueryRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

type BlocklistQueryResponse struct {
	giner.BasicResponse
	Data []g79.RealmsServerBlocklistPlayer `json:"data"` // 黑名单玩家列表
}

// BlocklistQuery godoc
//
//	@Summary		查询山头服黑名单
//	@Description	通过山头服 SID 查询黑名单玩家列表
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		BlocklistQueryRequest	true	"山头服 SID"
//	@Success		200				{object}	BlocklistQueryResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/blocklist_query [post]
func (*RealmsServer) BlocklistQuery(c *gin.Context) {
	var req BlocklistQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	players, ginerr := g79.QueryRealmsServerBlocklist(c.Request.Context(), gu, req.ServerID)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(players))
	c.Set("log", "查询山头服黑名单成功")
}
