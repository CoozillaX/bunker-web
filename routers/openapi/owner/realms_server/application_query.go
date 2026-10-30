package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApplicationQueryRequest struct {
	ServerID string `json:"server_id" binding:"min=1"` // 山头服 SID
}

type ApplicationQueryResponse struct {
	giner.BasicResponse
	Data map[string]any `json:"data"` // 申请队列, 主要字段为 queue
}

// ApplicationQuery godoc
//
//	@Summary		查询山头服加入申请
//	@Description	获取服主待审核的加入申请队列
//	@Description	队列项包含 application_id、applicant_uid、source、source_card_inviter、submit_time
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		ApplicationQueryRequest	true	"山头服 SID"
//	@Success		200				{object}	ApplicationQueryResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/application_query [post]
func (*RealmsServer) ApplicationQuery(c *gin.Context) {
	var req ApplicationQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	result, ginerr := g79.QueryRealmsServerApplications(c.Request.Context(), gu, req.ServerID)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(result))
	c.Set("log", "查询山头服加入申请成功")
}
