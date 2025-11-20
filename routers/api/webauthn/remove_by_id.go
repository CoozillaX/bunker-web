package webauthn

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/webauthn_credential"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RemoveRequest struct {
	CredentialID uint `json:"credential_id" binding:"gt=0"`
}

func (w *Webauthn) Remove(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Parse request
	var req RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check credential user id if match
	credential, err := webauthn_credential.QueryModelByID(req.CredentialID)
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	if credential.UserID != usr.ID {
		c.Error(giner.NewPublicGinError("删除失败"))
		return
	}
	// Remove
	if err := webauthn_credential.Remove(credential); err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("删除成功"))
	// Create log
	c.Set("log", "成功删除 Webauthn")
}
