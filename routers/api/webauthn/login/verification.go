package login

import (
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/webauthner"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *Login) Verification(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	// Begin discoverable login
	usr, ginerr := webauthner.FinishDiscoverableLogin(bearer.(string), c.Request)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("Welcome! "+usr.Username))
	// Set session
	session.Store("isPhoenix", false)
	session.Store("usr", usr)
	sessions.BindSessionToUsername(bearer.(string), usr.Username)
	// Create log
	c.Set("log", "Webauthn 登录成功")
}
