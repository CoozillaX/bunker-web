package phoenix

import (
	"bunker-core/protocol/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*Phoenix) TransferStartType(c *gin.Context) {
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	// Check if is phoenix
	if isPhoenix, ok := session.Load("isPhoenix"); !ok || !isPhoenix.(bool) {
		c.Error(giner.NewPublicGinError("无效会话"))
		return
	}
	// Get entityID
	entityID, ok := session.Load("entityID")
	if !ok {
		c.Error(giner.NewPublicGinError("会话已失效, 请重新登录"))
		return
	}
	// Get start type
	result, err := g79.GetStartType(entityID.(string), c.Query("content"))
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Return result
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(result))
	// Create log
	c.Set("log", "StartType 获取成功")
}
