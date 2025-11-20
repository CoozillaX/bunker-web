package login

import (
	"bunker-web/pkg/giner"
	"bunker-web/pkg/webauthner"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *Login) Options(c *gin.Context) {
	// Get bearer
	bearer, _ := c.Get("bearer")
	// Begin discoverable login
	option, err := webauthner.BeginDiscoverableLogin(bearer.(string))
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(option))
	// Create log
	c.Set("log", "请求 Webauthn 登录")
}
