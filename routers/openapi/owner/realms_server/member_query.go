package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberQueryRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

type MemberQueryResponse struct {
	giner.BasicResponse
	Data []map[string]any `json:"data"` // 成员列表
}

// MemberQuery godoc
//
//	@Summary		查询山头服成员
//	@Description	通过山头服 SID 查询成员列表
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		MemberQueryRequest	true	"山头服 SID"
//	@Success		200				{object}	MemberQueryResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/member_query [post]
func (*RealmsServer) MemberQuery(c *gin.Context) {
	var req MemberQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	members, ginerr := g79.QueryRealmsServerMembers(gu, req.ServerID)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(members))
	c.Set("log", "查询山头服成员成功")
}
