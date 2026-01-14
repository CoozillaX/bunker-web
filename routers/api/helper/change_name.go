package helper

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChangeNameRequest struct {
	UserName string `json:"username" binding:"min=1"`
}

func (*Helper) ChangeName(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req ChangeNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check helper if exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("未创建辅助用户"))
		return
	}
	// Store to DB
	defer models.DBSave(usr.HelperMpayUser)
	// Get g79 user
	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Change name
	if ginerr := g79.ChangeUserName(gu, req.UserName); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("已成功修改辅助用户昵称"))
	// Create log
	c.Set("log", fmt.Sprintf("更换Helper昵称成功, 新昵称: %s", req.UserName))
}
