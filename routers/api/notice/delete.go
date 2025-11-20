package notice

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/announcement"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	ID uint `json:"id" binding:"gt=0"`
}

func (*Notice) Delete(c *gin.Context) {
	// Parse request
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Delete notice
	if ginerr := announcement.DeleteByID(req.ID); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("公告删除成功"))
	// Create log
	c.Set("log", fmt.Sprintf("公告发布成功, id: %v", req.ID))
}
