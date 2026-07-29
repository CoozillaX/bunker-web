package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberSetPermissionRequest struct {
	ServerID   string `json:"server_id" binding:"min=1"`                             // 山头服 SID
	UID        string `json:"uid" binding:"min=1"`                                   // 玩家 UID
	Permission string `json:"permission" binding:"oneof=11111111 11111100 00000000"` // 权限: 管理员、成员、访客
}

// MemberSetPermission godoc
//
//	@Summary		设置山头服成员权限
//	@Description	权限值: 11111111=管理员, 11111100=成员, 00000000=访客
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		MemberSetPermissionRequest	true	"成员权限参数"
//	@Success		200				{object}	giner.BasicResponse			"成功时返回"
//	@Router			/openapi/owner/realms_server/member_set_permission [post]
func (*RealmsServer) MemberSetPermission(c *gin.Context) {
	var req MemberSetPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.SetRealmsServerMemberPermission(gu, req.ServerID, req.UID, req.Permission); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "设置山头服成员权限成功")
}
