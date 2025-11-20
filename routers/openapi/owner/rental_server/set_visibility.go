package rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetVisibilityRequest struct {
	ServerCode string `json:"server_code" binding:"min=1"`      // 服务器号
	Visibility int    `json:"visibility" binding:"gte=0,lte=2"` // 要设置的公开性, 0为所有人可见, 1为好友可见, 2为私密
}

// SetVisibility godoc
//
//	@Summary		设置服务器公开性
//	@Description	设置自己服务器的公开性, 需要提供 API Key 且绑定游戏账号
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			租赁服 (管理类)
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		SetVisibilityRequest	true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	giner.BasicResponse	"成功时返回"
//	@Router			/openapi/owner/rental_server/set_visibility [post]
func (*RentalServer) SetVisibility(c *gin.Context) {
	// Parse request
	var req SetVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	// Check owner if exists
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.MpayToken == "" {
		c.Error(giner.NewPublicGinError("请先绑定游戏账号"))
		return
	}
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser.MpayUser, nil)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Try to set level limitation
	if ginerr := g79.SetRentalServerVisibility(gu, req.ServerCode, req.Visibility); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	// Create log
	c.Set("log", "服务器公开性设置成功")
}
