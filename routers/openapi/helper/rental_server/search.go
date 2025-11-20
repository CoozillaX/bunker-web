package rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	ServerName string `json:"server_name" binding:"min=1"` // 要查找的租赁服号
}

type SearchResponse struct {
	giner.BasicResponse
	Data *g79.RentalServerInfo `json:"data"`
}

// SearchRentalServer godoc
//
//	@Summary		查询租赁服信息
//	@Description	查询租赁服基础信息, 需要提供 API key
//	@Description	注意: 此 API 会尝试进行游戏登录 (Helper)
//	@Tags			租赁服 (查询类)
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		SearchRequest	true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	SearchResponse	"成功时返回"
//	@Router			/openapi/helper/rental_server/search [post]
func (*RentalServer) Search(c *gin.Context) {
	// Parse request
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.MpayToken == "" {
		c.Error(giner.NewPublicGinError("未创建辅助用户"))
		return
	}
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser.MpayUser, nil)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Get player info by name
	serverInfo, ginerr := g79.QueryRentalServer(gu, req.ServerName)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Pack player info
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(serverInfo))
	// Create log
	c.Set("log", "查询租赁服成功")
}
