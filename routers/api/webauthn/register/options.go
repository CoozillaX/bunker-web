package register

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/webauthner"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *Register) Options(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Begin registration
	option, err := webauthner.BeginRegistration(bearer.(string), usr)
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(option))
	// Create log
	c.Set("log", "请求注册 Webauthn")
}
