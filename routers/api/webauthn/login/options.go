package login

import (
	"bunker-web/configs"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/webauthner"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (w *Login) Options(c *gin.Context) {
	// Begin discoverable login
	reqId, option, err := webauthner.BeginDiscoverableLogin()
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Set reqId cookie
	domain, _ := url.Parse(configs.CURRENT_WEB_DOMAIN)
	c.SetCookie(
		webauthner.WEBAUTHN_REQID_COOKIE_NAME,
		reqId,
		int(webauthner.WEBAUTHN_REQ_EXPIRE_TIME.Seconds()),
		"/",
		domain.Hostname(),
		domain.Scheme == "https",
		true,
	)
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(option))
	// Create log
	c.Set("log", "请求 Webauthn 登录")
}
