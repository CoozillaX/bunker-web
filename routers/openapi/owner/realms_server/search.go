package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	ServerCode string `json:"server_code" binding:"min=2" example:"#12345678"` // 游戏内显示的山头码, 需要包含 # 前缀
}

type SearchResponse struct {
	giner.BasicResponse
	Data map[string]any `json:"data"` // 山头服详细信息
}

// Search godoc
//
//	@Summary		查询山头服
//	@Description	通过游戏内带 # 前缀的山头码查询山头服详细信息
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		SearchRequest	true	"需要查询的山头码, 例如 #12345678"
//	@Success		200				{object}	SearchResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/search [post]
func (*RealmsServer) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	if !strings.HasPrefix(req.ServerCode, "#") || len(req.ServerCode) == 1 {
		c.Error(giner.NewPublicGinError("无效的山头码, 请包含 # 前缀"))
		return
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	result, ginerr := g79.QueryRealmsServer(gu, req.ServerCode)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(result))
	c.Set("log", "查询山头服详情成功")
}
