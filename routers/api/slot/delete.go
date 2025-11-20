package slot

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/slot"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	ID uint `json:"id" binding:"gt=0"`
}

func (*Slot) Delete(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Query slot
	s, err := slot.QueryByID(req.ID)
	if err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check permission
	if s.UserID != usr.ID {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Delete slot
	if ginerr := slot.Delete(s); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("删除成功"))
	// Create log
	c.Set("log", "删除 slot 成功")
}
