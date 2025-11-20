package user

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExtendExpireTimeRequest struct {
	UserName string `json:"username" binding:"min=1"`
	Seconds  int64  `json:"seconds" binding:"ne=0"`
}

func (*User) ExtendExpireTime(c *gin.Context) {
	// Parse request
	var req ExtendExpireTimeRequest
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
	// Renew user
	if ginerr := user.ExtendExpireTime(usr, req.Seconds); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("设置用户有效期成功"))
	// Create log
	c.Set("log", fmt.Sprintf("管理权限设置用户有效期成功, 目标用户名: %s, 时间: %v", req.UserName, req.Seconds))
}
