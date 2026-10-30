package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InviteCodeQueryRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

type InviteCodeQueryResponse struct {
	giner.BasicResponse
	Data []g79.RealmsServerInviteCode `json:"data"` // 邀请码列表
}

// InviteCodeQuery godoc
//
//	@Summary		查询山头服邀请码
//	@Description	通过山头服 SID 查询邀请码列表
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		InviteCodeQueryRequest	true	"山头服 SID"
//	@Success		200				{object}	InviteCodeQueryResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/invite_code_query [post]
func (*RealmsServer) InviteCodeQuery(c *gin.Context) {
	var req InviteCodeQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	codes, ginerr := g79.QueryRealmsServerInviteCodes(c.Request.Context(), gu, req.ServerID)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(codes))
	c.Set("log", "查询山头服邀请码成功")
}
