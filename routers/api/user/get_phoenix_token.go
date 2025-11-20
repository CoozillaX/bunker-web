package user

import (
	"bunker-web/models"
	"bunker-web/pkg/fbtoken"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPhoenixTokenRequest struct {
	HashedIP string `json:"hashed_ip" binding:"omitempty,len=32"`
}

func (*User) GetPhoenixToken(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req GetPhoenixTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get fbtoken
	token, _ := fbtoken.Encrypt(usr.Username, usr.Password, req.HashedIP)
	// Return as file
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", "attachment; filename=fbtoken")
	c.String(http.StatusOK, token)
	// Create log
	c.Set("log", "获取fbtoken")
}
