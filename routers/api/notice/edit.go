package notice

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/services/announcement"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EditRequest struct {
	ID      int    `json:"id" binding:"gt=0"`
	Title   string `json:"title" binding:"min=1"`
	Content string `json:"content" binding:"min=1"`
}

func (*Notice) Edit(c *gin.Context) {
	// Parse request
	var req EditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Query notice
	notice, ginerr := announcement.QueryByID(uint(req.ID))
	if ginerr != nil {
		c.Error(giner.NewPublicGinError("无效的公告ID"))
		return
	}
	// Modify notice info
	notice.Title = req.Title
	notice.Content = req.Content
	// Update to DB
	if ginerr := models.DBSave(notice); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("公告更新成功"))
	// Create log
	c.Set("log", fmt.Sprintf("公告更新成功, title: %s", req.Title))
}
