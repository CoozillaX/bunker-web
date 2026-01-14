package rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPlayerListRequest struct {
	Status    *int   `json:"status"`                    // 玩家状态 (是否黑名单) (0, 1, null(不限))
	ServerID  string `json:"server_id" binding:"min=1"` // 租赁服实体ID
	OrderType int    `json:"order_type"`                // 排序类型 (0, 1)
	IsOnline  *bool  `json:"is_online"`                 // 是否在线 (false, true, null(不限))
	Length    *int   `json:"length" example:"50"`       // 查询分页长度 (null(不限))
	Offset    int    `json:"offset"`                    // 查询分页偏移
}

type GetPlayerListResponse struct {
	giner.BasicResponse
	Data []map[string]any `json:"data"`
}

// GetRentalServerPlayerList godoc
//
//	@Summary		获取租赁服玩家列表
//	@Description	获取租赁服玩家列表, 例如历史进入列表, 需要提供 API key
//	@Description	注意: 此 API 会尝试进行游戏登录 (Helper)
//	@Tags			租赁服 (查询类)
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		GetPlayerListRequest	true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	GetPlayerListResponse	"成功时返回"
//	@Router			/openapi/helper/rental_server/get_player_list [post]
func (*RentalServer) GetPlayerList(c *gin.Context) {
	// Parse request
	var req GetPlayerListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("未创建辅助用户"))
		return
	}
	// g79 login
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Get player info by name
	result, ginerr := g79.QueryRentalServerPlayerList(gu, req.Status, req.ServerID, req.OrderType, req.IsOnline, req.Length, req.Offset)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Pack player info
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(result))
	// Create log
	c.Set("log", "查询租赁服玩家列表成功")
}
