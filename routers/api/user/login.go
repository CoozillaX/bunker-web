package user

import (
	"bunker-web/configs"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/pkg/utils"
	"bunker-web/services/user"
	"bunker-web/services/user_ban_record"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	UserName string `json:"username" binding:"min=1"`
	Password string `json:"password" binding:"len=64"`
}

func (*User) Login(c *gin.Context) {
	// Parse request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Login
	usr, ginerr := user.NormalLogin(req.UserName, utils.SHA256Hex([]byte(req.Password+configs.USER_PSW_SALT)))
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Ban Check
	if banRecord, _ := user_ban_record.GetCurrentBanRecordFormattedStringByUserID(usr.ID); len(banRecord) > 0 {
		c.Error(giner.NewPublicGinError(banRecord))
		return
	}
	// Create bearer
	bearer := uuid.NewString()
	// Create session
	session := sessions.CreateSessionByBearer(bearer)
	sessions.BindSessionToUsername(bearer, usr.Username)
	session.Store("isPhoenix", false)
	session.Store("usr", usr)
	// Set cookie
	domain, _ := url.Parse(configs.CURRENT_WEB_DOMAIN)
	c.SetCookie(
		sessions.SESSION_COOKIE_NAME,
		bearer,
		int(sessions.SESSION_EXPIRE_TIME.Seconds()),
		"/",
		domain.Hostname(),
		domain.Scheme == "https",
		true,
	)
	// Response
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("Welcome! "+usr.Username))
	// Create log
	c.Set("log", "用户中心登录成功")
}
