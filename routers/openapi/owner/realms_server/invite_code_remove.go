package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InviteCodeRemoveRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
	Code     string `json:"code" binding:"min=1"`      // 邀请码
}

// InviteCodeRemove godoc
//
//	@Summary		删除山头服邀请码
//	@Description	删除指定山头服的邀请码
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		InviteCodeRemoveRequest	true	"山头服 SID 和邀请码"
//	@Success		200				{object}	giner.BasicResponse		"成功时返回"
//	@Router			/openapi/owner/realms_server/invite_code_remove [post]
func (*RealmsServer) InviteCodeRemove(c *gin.Context) {
	var req InviteCodeRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.RemoveRealmsServerInviteCode(gu, req.ServerID, req.Code); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "删除山头服邀请码成功")
}
