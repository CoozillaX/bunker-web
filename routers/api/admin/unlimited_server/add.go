package unlimited_server

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/unlimited_rental_server"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	ServerCode string `json:"server_code" binding:"min=1"`
}

func (*UnlimitedServer) Add(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check if exists
	if _, ginerr := unlimited_rental_server.QueryByServerCode(req.ServerCode); ginerr == nil {
		c.Error(giner.NewPublicGinError("此服务器已被设置为无限制进入, 无需重复设置"))
		return
	}
	// Add
	if _, ginerr := unlimited_rental_server.Create(usr.ID, req.ServerCode); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("设置成功"))
	// Create log
	c.Set("log", fmt.Sprintf("管理权限设置无限制服务器成功, 服务器号: %s", req.ServerCode))
}
