package user

import (
	"bunker-web/models"
	"bunker-web/pkg/fbtoken"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*User) GetPhoenixToken(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Get fbtoken
	token, _ := fbtoken.Encrypt(usr.Username, usr.Password)
	// Return as file
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=fbtoken")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.String(http.StatusOK, token)
	// Create log
	c.Set("log", "获取fbtoken")
}
