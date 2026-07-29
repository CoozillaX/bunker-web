package realms_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberRemoveRequest struct {
	ServerID  string `json:"server_id" binding:"min=1"` // 山头服 SID
	UID       string `json:"uid"`                       // 玩家 UID, remove_all=false 时必填
	RemoveAll bool   `json:"remove_all"`                // 是否移除全部成员
}

// MemberRemove godoc
//
//	@Summary		移除山头服成员
//	@Description	移除指定成员; remove_all=true 时移除全部成员且无需传 uid
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			山头服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		MemberRemoveRequest	true	"成员移除参数"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/realms_server/member_remove [post]
func (*RealmsServer) MemberRemove(c *gin.Context) {
	var req MemberRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil || (!req.RemoveAll && req.UID == "") {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	if ginerr = g79.RemoveRealmsServerMember(gu, req.ServerID, req.UID, req.RemoveAll); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	c.Set("log", "移除山头服成员成功")
}
