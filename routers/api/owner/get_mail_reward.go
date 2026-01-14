package owner

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*Owner) GetMailReward(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	u, _ := session.Load("usr")
	usr := u.(*models.User)
	// Check owner if exist
	if usr.OwnerMpayUser == nil || usr.OwnerMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("未创建游戏账号"))
		return
	}
	// Store to DB
	defer models.DBSave(usr.OwnerMpayUser)
	// Relogin
	gu, ginerr := g79.HandleG79Login(usr.OwnerMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Get all mails
	mailList, ginerr := g79.GetGeneralMailList(gu)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Read all mails
	if ginerr := g79.ReadMails(gu, mailList); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Claim all rewards
	if ginerr := g79.ClaimMailRewards(gu, mailList); ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("已领取所有邮件奖励"))
	// Create log
	c.Set("log", "成功领取邮件奖励")
}
