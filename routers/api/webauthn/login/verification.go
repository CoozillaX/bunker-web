package login

import (
	"bunker-web/configs"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/webauthner"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (w *Login) Verification(c *gin.Context) {
	// Get cookie
	reqId, err := c.Cookie(webauthner.WEBAUTHN_REQID_COOKIE_NAME)
	if err != nil {
		c.Error(giner.NewPublicGinError("无效请求"))
		return
	}
	// Begin discoverable login
	usr, ginerr := webauthner.FinishDiscoverableLogin(reqId, c.Request)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Create bearer
	bearer := uuid.NewString()
	// Set session
	session := sessions.CreateSessionByBearer(bearer)
	sessions.BindSessionToUsername(bearer, usr.Username)
	session.Store("isPhoenix", false)
	session.Store("usr", usr)
	// Clear reqId cookie
	domain, _ := url.Parse(configs.CURRENT_WEB_DOMAIN)
	c.SetCookie(
		webauthner.WEBAUTHN_REQID_COOKIE_NAME,
		reqId,
		-1,
		"/",
		domain.Hostname(),
		domain.Scheme == "https",
		true,
	)
	// Set session cookie
	c.SetCookie(
		sessions.SESSION_COOKIE_NAME,
		bearer,
		int(sessions.SESSION_EXPIRE_TIME.Seconds()),
		"/",
		domain.Hostname(),
		domain.Scheme == "https",
		true,
	)
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("Welcome! "+usr.Username))
	// Create log
	c.Set("log", "Webauthn 登录成功")
}
