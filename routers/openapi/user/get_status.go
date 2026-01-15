package user

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Slot struct {
	GameID   int    `json:"game_id"`   // Slots 绑定的游戏ID
	ExpireAt int64  `json:"expire_at"` // Slots 有效期至
	Note     string `json:"note"`      // Slots 备注
}

type GetStatusResponseData struct {
	Username string  `json:"username"`  // 用户名
	GameID   int     `json:"game_id"`   // 绑定的游戏ID
	IsAdmin  bool    `json:"is_admin"`  // 是否为系统管理员
	CreateAt int64   `json:"create_at"` // 创建时间
	Slots    []*Slot `json:"slots"`     // Slots 列表
}

type GetStatusResponse struct {
	giner.BasicResponse
	Data *GetStatusResponseData `json:"data"`
}

// GetStatus godoc
//
//	@Summary		用户信息查询
//	@Description	用户信息查询, 需要提供 API Key
//	@Tags			用户中心
//	@Accept			json
//	@Produce		json
//
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//
//	@Success		200				{object}	GetStatusResponse	"成功时返回"
//	@Router			/openapi/user/get_status [get]
func (*User) GetStatus(c *gin.Context) {
	// Get user
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	// Pack user info
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(&GetStatusResponseData{
		Username: usr.Username,
		GameID:   usr.GameID,
		IsAdmin:  usr.Permission == user.PermissionAdmin,
		CreateAt: usr.CreatedAt.UnixMilli(),
	}))
	// Create log
	c.Set("log", "获取usr信息成功")
}
