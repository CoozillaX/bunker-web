package log

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	go func() {
		for {
			DeleteOldLogs()
			time.Sleep(time.Hour * 24)
		}
	}()
}

func DeleteOldLogs() error {
	expiredDate := time.Now().AddDate(0, -1, 0)
	return models.DB.Where("created_at < ?", expiredDate).Delete(&models.Log{}).Error
}

func GinLogCreate(c *gin.Context) {
	// 1. Try to get user
	var usr *models.User
	{
		if bearer, ok := c.Get("bearer"); ok {
			if session, ok := sessions.GetSessionByBearer(bearer.(string)); ok {
				if u, ok := session.Load("usr"); ok {
					usr = u.(*models.User)
				}
			}
		}
	}
	// 2. Try to get gin error
	var publicError, privateError string
	{
		if ginerr := c.Errors.Last(); ginerr != nil {
			publicError = giner.GetPublicErrorString(ginerr)
			privateError = giner.GetPrivateErrorString(ginerr)
		}
	}
	// 3. Try to get extra message
	var usrID uint
	var message string
	{
		if usr != nil {
			usrID = usr.ID
			message += fmt.Sprintf("用户名(%s) ", usr.Username)
		}
		message += c.GetString("log")
	}
	// 4. Create log
	models.DBCreate(&models.Log{
		IPAddress:    c.ClientIP(),
		Method:       c.Request.Method,
		Path:         c.Request.URL.Path,
		UserID:       usrID,
		PublicError:  publicError,
		PrivateError: privateError,
		Message:      message,
	})
}
