package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApplicationApproveItem struct {
	ApplicationID int64   `json:"application_id" binding:"min=1"` // 申请 ID
	Permission    *string `json:"permission,omitempty"`           // 通过时的成员权限, 不传则使用默认权限
}

type ApplicationApproveRequest struct {
	ServerID string                   `json:"server_id" binding:"min=1"`                  // 山头服 SID
	Action   int                      `json:"action" binding:"oneof=1 2 3"`               // 1=通过, 2=拒绝, 3=拉黑并拒绝
	Items    []ApplicationApproveItem `json:"items" binding:"required,min=1,max=20,dive"` // 待审批项目, 最多 20 项
}

type ApplicationApproveResponse struct {
	giner.BasicResponse
	Data map[string]any `json:"data"` // 审批结果, 主要字段为 results
}

// ApplicationApprove godoc
//
//	@Summary		审批山头服加入申请
//	@Description	批量审批加入申请, 每次最多 20 项
//	@Description	action: 1=通过, 2=拒绝, 3=拉黑并拒绝; permission 仅 action=1 时生效
//	@Description	permission: 11111111=管理员, 11111100=成员, 00000000=访客
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		ApplicationApproveRequest	true	"批量审批参数"
//	@Success		200				{object}	ApplicationApproveResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/application_approve [post]
func (*RealmsServer) ApplicationApprove(c *gin.Context) {
	var req ApplicationApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil || !validApplicationPermissions(req.Action, req.Items) {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	items := make([]map[string]any, 0, len(req.Items))
	for _, item := range req.Items {
		value := map[string]any{"application_id": item.ApplicationID}
		if req.Action == 1 && item.Permission != nil {
			value["permission"] = *item.Permission
		}
		items = append(items, value)
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	result, ginerr := g79.ApproveRealmsServerApplications(c.Request.Context(), gu, req.ServerID, req.Action, items)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(result))
	c.Set("log", "审批山头服加入申请成功")
}

func validApplicationPermissions(action int, items []ApplicationApproveItem) bool {
	for _, item := range items {
		if item.Permission == nil {
			continue
		}
		if action != 1 {
			return false
		}
		switch *item.Permission {
		case "11111111", "11111100", "00000000":
		default:
			return false
		}
	}
	return true
}
