package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InviteCodeGenerateRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

type InviteCodeGenerateResponse struct {
	giner.BasicResponse
	Data string `json:"data"` // 新生成的邀请码
}

// InviteCodeGenerate godoc
//
//	@Summary		生成山头服邀请码
//	@Description	为指定山头服生成永久有效的邀请码
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		InviteCodeGenerateRequest	true	"山头服 SID"
//	@Success		200				{object}	InviteCodeGenerateResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/invite_code_generate [post]
func (*RealmsServer) InviteCodeGenerate(c *gin.Context) {
	var req InviteCodeGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	code, ginerr := g79.GenerateRealmsServerInviteCode(gu, req.ServerID)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(code))
	c.Set("log", "生成山头服邀请码成功")
}
