package user

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetPermissionRequest struct {
	UserName   string `json:"username" binding:"min=1"`
	Permission uint   `json:"permission" binding:"gte=0,lt=2"`
}

func (*User) SetPermission(c *gin.Context) {
	// Parse request
	var req SetPermissionRequest
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
	// Set and save user permission
	usr.Permission = req.Permission
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("更改用户权限成功"))
	// Create log
	c.Set("log", fmt.Sprintf("管理权限更改用户权限成功, 目标用户名: %s, 新权限: %v", req.UserName, req.Permission))
}
