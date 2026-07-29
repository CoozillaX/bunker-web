package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberInviteRequest struct {
	ServerID   string `json:"server_id" binding:"min=1"`                             // 山头服 SID
	UID        string `json:"uid" binding:"min=1"`                                   // 被邀请玩家 UID
	Permission string `json:"permission" binding:"oneof=11111111 11111100 00000000"` // 加入后的权限
}

// MemberInvite godoc
//
//	@Summary		邀请山头服成员
//	@Description	服主主动邀请指定玩家加入山头服
//	@Description	permission: 11111111=管理员, 11111100=成员, 00000000=访客
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		MemberInviteRequest	true	"成员邀请参数"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/member_invite [post]
func (*RealmsServer) MemberInvite(c *gin.Context) {
	var req MemberInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.InviteRealmsServerMember(gu, req.ServerID, req.UID, req.Permission); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "邀请山头服成员成功")
}
