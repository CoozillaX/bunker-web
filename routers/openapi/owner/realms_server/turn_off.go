package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TurnOffRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

// TurnOff godoc
//
//	@Summary		关闭山头服
//	@Description	通过山头服 SID 关闭指定山头服
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		TurnOffRequest	true	"山头服 SID"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/turn_off [post]
func (*RealmsServer) TurnOff(c *gin.Context) {
	var req TurnOffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.TurnOffRealmsServer(c.Request.Context(), gu, req.ServerID); ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "请求关闭山头服成功")
}
