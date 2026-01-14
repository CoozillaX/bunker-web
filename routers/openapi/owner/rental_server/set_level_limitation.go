package rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetLevelLimitationRequest struct {
	ServerCode string `json:"server_code" binding:"min=1"`  // 服务器号
	Level      int    `json:"level" binding:"gte=0,lte=50"` // 要设置的服务器准入等级, 0为关闭, 最高为50(基岩V)
}

// SetLevelLimitation godoc
//
//	@Summary		设置服务器等级限制
//	@Description	设置自己服务器的等级限制, 需要提供 API Key 且绑定游戏账号
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			租赁服 (管理类)
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		SetLevelLimitationRequest	true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	giner.BasicResponse		"成功时返回"
//	@Router			/openapi/owner/rental_server/set_level_limitation [post]
func (*RentalServer) SetLevelLimitation(c *gin.Context) {
	// Parse request
	var req SetLevelLimitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	// Check owner if exists
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("请先绑定游戏账号"))
		return
	}
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Try to set level limitation
	if ginerr := g79.SetRentalServerLevelLimitation(gu, req.ServerCode, req.Level); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true))
	// Create log
	c.Set("log", "服务器等级限制设置成功")
}
