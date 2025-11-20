package register

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/webauthner"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *Register) Verification(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Finish registration
	if ginerr := webauthner.FinishRegistration(bearer.(string), usr, c.Request); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("添加成功"))
	// Create log
	c.Set("log", "成功注册 Webauthn")
}
