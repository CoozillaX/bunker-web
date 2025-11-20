package middlewares

import (
	"bunker-web/models"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"

	"bunker-web/pkg/giner"

	"github.com/gin-gonic/gin"
)

func NormalPermissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session
		bearer, _ := c.Get("bearer")
		session, _ := sessions.GetSessionByBearer(bearer.(string))
		u, _ := session.Load("usr")
		usr := u.(*models.User)
		// Check normal account status
		if vaild, reason := user.CheckIfVaild(usr); !vaild {
			c.AbortWithStatusJSON(http.StatusForbidden, giner.MakeHTTPResponse(false).SetMessage(reason))
			return
		}
		c.Next()
	}
}

func AdminPermissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session
		bearer, _ := c.Get("bearer")
		session, _ := sessions.GetSessionByBearer(bearer.(string))
		u, _ := session.Load("usr")
		usr := u.(*models.User)
		// Check permission
		if usr.Permission != user.PermissionAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, giner.MakeHTTPResponse(false).SetMessage("无权访问"))
			return
		}
		c.Next()
	}
}
