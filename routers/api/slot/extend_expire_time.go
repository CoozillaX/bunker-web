package slot

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/redeem"
	"bunker-web/services/slot"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExtendExpireTimeRequest struct {
	ID         uint   `json:"id" binding:"gt=0"`
	RedeemCode string `json:"redeem_code" binding:"len=36"`
}

func (*Slot) ExtendExpireTime(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req ExtendExpireTimeRequest
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
	// Check if binded game id
	if s.GameID == 0 {
		c.Error(giner.NewPublicGinError("请先为 slot 绑定游戏ID"))
		return
	}
	// Redeem
	result, ginerr := redeem.UseRedeemCodeForSlot(usr, s, req.RedeemCode)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage(result))
	// Create log
	c.Set("log", result)
}
