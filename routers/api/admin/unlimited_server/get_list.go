package unlimited_server

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/unlimited_rental_server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnlimitedRentalServerInfo struct {
	ID         uint   `json:"id"`
	CreateAt   int64  `json:"create_at"`
	ServerCode string `json:"server_code"`
}

func (*UnlimitedServer) GetList(c *gin.Context) {
	// Get all
	svrs, ginerr := unlimited_rental_server.QueryAll()
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Format
	UnlimitedRentalServerInfoList := make([]UnlimitedRentalServerInfo, len(svrs))
	for i, svr := range svrs {
		UnlimitedRentalServerInfoList[i] = UnlimitedRentalServerInfo{
			ID:         svr.ID,
			CreateAt:   svr.CreatedAt.UnixMilli(),
			ServerCode: svr.ServerCode,
		}
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(UnlimitedRentalServerInfoList))
	// Create log
	c.Set("log", "管理权限获取无限制服务器列表成功")
}
