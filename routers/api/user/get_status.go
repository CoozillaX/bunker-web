package user

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"bunker-web/services/webauthn_credential"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Slot struct {
	ID       uint   `json:"id"`
	GameID   int    `json:"game_id"`
	ExpireAt int64  `json:"expire_at"`
	Note     string `json:"note"`
}

type CredentialInfo struct {
	ID       uint   `json:"id"`
	CreateAt int64  `json:"create_at"`
	RawID    string `json:"raw_id"`
}

type GetStatusResponseData struct {
	Username       string            `json:"username"`
	GameID         int               `json:"game_id"`
	Permission     uint              `json:"permission"`
	IsAdmin        bool              `json:"is_admin"`
	CreateAt       int64             `json:"create_at"`
	APIKey         string            `json:"api_key"`
	HasEmail       bool              `json:"has_email"`
	ClientUsername string            `json:"client_username"`
	Credentials    []*CredentialInfo `json:"credentials"`
}

func (*User) GetStatus(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Refresh user
	usr, _ = user.QueryByUsername(usr.Username)
	if usr == nil {
		c.Error(giner.NewPublicGinError("用户信息异常, 请重新登录"))
		return
	}
	session.Store("usr", usr)
	// Query Credentials
	credentials, err := webauthn_credential.QueryModelsByUserID(usr.ID)
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	credentialInfos := make([]*CredentialInfo, len(credentials))
	for i, credential := range credentials {
		credentialInfos[i] = &CredentialInfo{
			ID:       credential.ID,
			CreateAt: credential.CreatedAt.UnixMilli(),
			RawID:    credential.RawID,
		}
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(
		&GetStatusResponseData{
			Username:    usr.Username,
			GameID:      usr.GameID,
			Permission:  usr.Permission,
			IsAdmin:     usr.Permission == user.PermissionAdmin,
			CreateAt:    usr.CreatedAt.UnixMilli(),
			APIKey:      usr.APIKey,
			HasEmail:    usr.Email != "",
			Credentials: credentialInfos,
		},
	))
	// Create log
	c.Set("log", "获取usr信息成功")
}
