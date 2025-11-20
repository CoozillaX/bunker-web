package unlimited_server

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/unlimited_rental_server"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	RentalServerID uint `json:"rental_server_id" binding:"gt=0"`
}

func (*UnlimitedServer) Delete(c *gin.Context) {
	// Parse request
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check if exists
	svr, ginerr := unlimited_rental_server.QueryByID(req.RentalServerID)
	if ginerr != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Delete
	if ginerr := unlimited_rental_server.Delete(svr); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("删除成功"))
	// Create log
	c.Set("log", fmt.Sprintf("管理权限删除无限制服务器成功, 服务器号: %s", svr.ServerCode))
}
