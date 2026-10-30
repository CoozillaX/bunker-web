package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberInviteRemoveRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
	UID      string `json:"uid" binding:"min=1"`       // 被邀请玩家 UID
}

// MemberInviteRemove godoc
//
//	@Summary		撤回山头服成员邀请
//	@Description	撤回服主发出的待处理成员邀请
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		MemberInviteRemoveRequest	true	"成员邀请撤回参数"
//	@Success		200				{object}	giner.BasicResponse		"成功时返回"
//	@Router			/openapi/owner/realms_server/member_invite_remove [post]
func (*RealmsServer) MemberInviteRemove(c *gin.Context) {
	var req MemberInviteRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.RemoveRealmsServerMemberInvite(c.Request.Context(), gu, req.ServerID, req.UID); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "撤回山头服成员邀请成功")
}
