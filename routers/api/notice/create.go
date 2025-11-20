package notice

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/announcement"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Title   string `json:"title" binding:"min=1"`
	Content string `json:"content" binding:"min=1"`
}

func (*Notice) Create(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Create announcement
	if ginerr := announcement.Create(usr, req.Title, req.Content); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("公告发布成功"))
	// Create log
	c.Set("log", fmt.Sprintf("公告发布成功, title: %s", req.Title))
}
