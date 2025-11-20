package user

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*User) Logout(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Always ok because of auth middleware
	sessions.DeleteSessionByBearer(bearer.(string))
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("See ya!"))
	// Create log
	c.Set("log", fmt.Sprintf("用户名(%s) 用户中心登出成功", usr.Username))
}
