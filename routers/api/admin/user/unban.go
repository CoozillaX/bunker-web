package user

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"bunker-web/services/user_ban_record"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnBanRequest struct {
	UserName string `json:"username" binding:"min=1"`
}

func (*User) UnBan(c *gin.Context) {
	// Parse request
	var req UnBanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Query user
	usr, ginerr := user.QueryByUsername(req.UserName)
	if ginerr != nil {
		c.Error(giner.NewPublicGinError("用户不存在"))
		return
	}
	// Check user status
	if usr.Permission == user.PermissionAdmin {
		c.Error(giner.NewPublicGinError("无法对管理员进行操作"))
		return
	}
	// Remove all ban record
	if ginerr := user_ban_record.RemoveAllBanRecordByUserID(usr.ID); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("解封用户成功"))
	// Create log
	c.Set("log", fmt.Sprintf("管理权限解封用户成功, 目标用户名: %s", req.UserName))
}
