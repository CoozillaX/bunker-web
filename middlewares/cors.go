package middlewares

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var allowedHosts = map[string]bool{
	"liliya233.uk": true,
	"localhost":    true,
}

func CORSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		u, err := url.Parse(origin)
		if err != nil {
			c.Next()
			return
		}

		if _, ok := allowedHosts[u.Hostname()]; !ok {
			c.Next()
			return
		}

		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Accept, Content-Type, Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
