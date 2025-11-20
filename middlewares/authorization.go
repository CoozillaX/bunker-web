package middlewares

import (
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BearerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session by bearer
		bearer := c.GetHeader("Authorization")
		bearer = strings.TrimPrefix(bearer, "Bearer ")
		_, ok := sessions.GetSessionByBearer(bearer)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, giner.MakeHTTPResponse(false).SetMessage("会话已失效, 请重新登录"))
			return
		}
		c.Set("bearer", bearer)
	}
}

func OpenAPIKeyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session by bearer
		apiKey := c.GetHeader("Authorization")
		// Check usr if exists
		usr, ginerr := user.QueryUserByAPIKey(apiKey)
		if ginerr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, giner.MakeHTTPResponse(false).SetMessage("无效的API Key"))
			return
		}
		// Check if vaild
		if !strings.HasPrefix(c.Request.URL.Path, "/openapi/user/") {
			if vaild, reason := user.CheckIfVaild(usr); !vaild {
				c.AbortWithStatusJSON(http.StatusUnauthorized, giner.MakeHTTPResponse(false).SetMessage(reason))
				return
			}
		}
		// Store usr
		c.Set("usr", usr)
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session by bearer
		bearer, _ := c.Get("bearer")
		session, _ := sessions.GetSessionByBearer(bearer.(string))
		// Check usr if exists
		if _, ok := session.Load("usr"); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, giner.MakeHTTPResponse(false).SetMessage("会话已失效, 请重新登录"))
			return
		}
		// Phoenix login allow visit Phoenix api only
		isPhoenix, ok := session.Load("isPhoenix")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, giner.MakeHTTPResponse(false).SetMessage("会话已失效, 请重新登录"))
			return
		}
		if isPhoenix.(bool) && !strings.HasPrefix(c.Request.URL.Path, "/api/phoenix/") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, giner.MakeHTTPResponse(false).SetMessage("会话已失效, 请重新登录"))
			return
		}
		c.Next()
	}
}
