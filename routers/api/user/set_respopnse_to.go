package user

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetResponseToRequest struct {
	ClientUsername string `json:"client_username" binding:"max=20"`
}

func (*User) SetResponseTo(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req SetResponseToRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Set response to
	usr.ResponseTo = req.ClientUsername
	// Save to db
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("设置成功"))
	// Create log
	c.Set("log", "设置 response_to 成功")
}
