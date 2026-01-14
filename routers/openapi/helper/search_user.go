package helper

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchUserRequest struct {
	Nickname string `json:"nickname" binding:"min=1"` // 要查找的玩家昵称
}

type SearchUserResponse struct {
	giner.BasicResponse
	Data []map[string]any `json:"data"`
}

// SearchUser godoc
//
//	@Summary		查询玩家
//	@Description	查询玩家, 可通过昵称来获取uid, 需要提供 API key
//	@Description	注意: 此 API 会尝试进行游戏登录 (Helper)
//	@Tags			玩家
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		SearchUserRequest		true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	SearchUserResponse	"成功时返回"
//	@Router			/openapi/helper/search_user [post]
func (*Helper) SearchUser(c *gin.Context) {
	// Parse request
	var req SearchUserRequest
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
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Get player info by name
	result, ginerr := g79.GetPlayerInfoByName(gu, req.Nickname)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Pack player info
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(result))
	// Create log
	c.Set("log", "查询玩家成功")
}
