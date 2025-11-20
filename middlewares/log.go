package middlewares

import (
	"bunker-web/services/log"

	"github.com/gin-gonic/gin"
)

func LogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.GinLogCreate(c)
	}
}
