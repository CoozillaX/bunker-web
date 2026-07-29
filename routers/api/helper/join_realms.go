package helper

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubmitRealmsApplicationRequest struct {
	Code string `json:"code" binding:"min=1,max=20"`
}

func (*Helper) SubmitRealmsApplication(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req SubmitRealmsApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	gu, ginerr := g79.GetHelperG79User(c)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Store to DB
	defer models.DBSave(usr.HelperMpayUser)
	// Submit realms server application
	if ginerr := gu.SubmitRealmsServerApplication(req.Code); ginerr != nil {
		c.Error(giner.NewPublicGinError(ginerr.Error()))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("已成功提交山头服申请"))
	// Create log
	c.Set("log", "提交山头服加入申请成功")
}
