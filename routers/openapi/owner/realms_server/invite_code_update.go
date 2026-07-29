package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InviteCodeUpdateRequest struct {
	ServerID  string  `json:"server_id" binding:"min=1"`     // 山头服 SID
	Code      string  `json:"code" binding:"min=1"`          // 邀请码
	ShareText *string `json:"share_text" binding:"required"` // 分享文案, 允许传空字符串
}

// InviteCodeUpdate godoc
//
//	@Summary		更新山头服邀请码
//	@Description	更新指定邀请码的分享文案
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		InviteCodeUpdateRequest	true	"山头服 SID、邀请码和分享文案"
//	@Success		200				{object}	giner.BasicResponse		"成功时返回"
//	@Router			/openapi/owner/realms_server/invite_code_update [post]
func (*RealmsServer) InviteCodeUpdate(c *gin.Context) {
	var req InviteCodeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.UpdateRealmsServerInviteCode(gu, req.ServerID, req.Code, *req.ShareText); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "更新山头服邀请码成功")
}
