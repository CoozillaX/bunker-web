package rental_server

import (
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryBackupRequest struct {
	ServerCode string `json:"server_code" binding:"min=1"` // 服务器号
}

type QueryBackupResponse struct {
	giner.BasicResponse
	Data []*g79.RentalServerBackupInfo `json:"data"` // 服务器备份列表
}

// QueryBackup godoc
//
//	@Summary		查询服务器备份
//	@Description	查询服务器备份及处理状态, 需要提供 API Key 且绑定游戏账号
//	@Description	注意: 此 API 会尝试进行游戏登录 (Owner)
//	@Tags			租赁服
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Param			Request			body		QueryBackupRequest	true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200				{object}	QueryBackupResponse		"成功时返回"
//	@Router			/openapi/owner/rental_server/query_backup [post]
func (*RentalServer) QueryBackup(c *gin.Context) {
	var req QueryBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}

	gu, ginerr := g79.GetOwnerG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	backupInfoList, ginerr := g79.QueryRentalServerBackup(c.Request.Context(), gu, req.ServerCode)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, &QueryBackupResponse{
		BasicResponse: giner.BasicResponse{Success: true, Message: "ok"},
		Data:          backupInfoList,
	})
	c.Set("log", "查询服务器备份成功")
}
