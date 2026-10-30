package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TurnOnRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

// TurnOn godoc
//
//	@Summary		开启山头服
//	@Description	通过山头服 SID 开启指定山头服
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		TurnOnRequest	true	"山头服 SID"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/turn_on [post]
func (*RealmsServer) TurnOn(c *gin.Context) {
	var req TurnOnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.TurnOnRealmsServer(c.Request.Context(), gu, req.ServerID); ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "请求开启山头服成功")
}
