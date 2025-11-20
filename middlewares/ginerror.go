package middlewares

import (
	"net/http"

	"bunker-web/pkg/giner"

	"github.com/gin-gonic/gin"
)

func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 0. Check body if already written
		if c.Writer.Written() {
			return
		}
		// 1. Init variables
		ginerr := c.Errors.Last()
		if ginerr == nil {
			return
		}
		// 2. Get public error string
		publicStr := giner.GetPublicErrorString(ginerr)
		if publicStr == "" {
			publicStr = "未知错误"
		}
		// 3. Response
		c.AbortWithStatusJSON(http.StatusOK, giner.MakeHTTPResponse(false).
			SetMessage(publicStr).
			SetTranslation(giner.GetTranslationCode(ginerr)),
		)
	}
}
