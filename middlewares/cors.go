package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin == "" {
			origin = "https://liliya233.uk"
		}

		if origin == "https://liliya233.uk" {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Accept, Content-Type, Cookie")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type")
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
