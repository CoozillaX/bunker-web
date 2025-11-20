package api_key

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateResponseData struct {
	APIKey string `json:"api_key"`
}

func (*APIKey) Generate(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Generate API key and ensure it is unique
	newKey := uuid.NewString()
	for {
		if user, _ := user.QueryUserByAPIKey(newKey); user == nil {
			break
		}
		newKey = uuid.NewString()
	}
	// Update user API key
	usr.APIKey = newKey
	if err := models.DBSave(usr); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(
		&GenerateResponseData{
			APIKey: newKey,
		},
	))
	// Create log
	c.Set("log", "获取API Key成功")
}
