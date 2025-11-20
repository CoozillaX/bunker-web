package api

import (
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (*API) New(c *gin.Context) {
	bearer := uuid.NewString()
	sessions.CreateSessionByBearer(bearer)
	c.String(http.StatusOK, bearer)
	// Create log
	c.Set("log", "创建Bearer成功")
}
