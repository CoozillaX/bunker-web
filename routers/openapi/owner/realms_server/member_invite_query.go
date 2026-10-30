package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberInviteQueryRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

type MemberInviteQueryResponse struct {
	giner.BasicResponse
	Data []map[string]any `json:"data"` // 待处理邀请列表
}

// MemberInviteQuery godoc
//
//	@Summary		查询山头服成员邀请
//	@Description	查询服主已发出且仍待处理的成员邀请
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		MemberInviteQueryRequest	true	"山头服 SID"
//	@Success		200				{object}	MemberInviteQueryResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/member_invite_query [post]
func (*RealmsServer) MemberInviteQuery(c *gin.Context) {
	var req MemberInviteQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	invites, ginerr := g79.QueryRealmsServerMemberInvites(c.Request.Context(), gu, req.ServerID)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(invites))
	c.Set("log", "查询山头服成员邀请成功")
}
