package user

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"bunker-web/services/user_ban_record"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	UserName string `json:"username" binding:"min=1"`
}

type QueryResponseData struct {
	UserName   string `json:"username"`
	GameID     int    `json:"game_id"`
	Permission uint   `json:"permission"`
	CreateAt   int64  `json:"create_at"`
	BanUntil   int64  `json:"ban_until"`
	BanReason  string `json:"ban_reason"`
}

func (*User) Query(c *gin.Context) {
	// Parse request
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Query user
	usr, ginerr := user.QueryByUsername(req.UserName)
	if ginerr != nil {
		c.Error(giner.NewPublicGinError("用户不存在"))
		return
	}
	// Query user ban record
	var banUntil int64
	var banReason string
	if usrBanRecord, _ := user_ban_record.GetCurrentBanRecordByUserID(usr.ID); usrBanRecord != nil {
		banUntil = usrBanRecord.Until.Time.UnixMilli()
		banReason = usrBanRecord.Reason
	}
	// Pack user info
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage("查询用户成功").SetData(&QueryResponseData{
		UserName:   usr.Username,
		GameID:     usr.GameID,
		Permission: usr.Permission,
		CreateAt:   usr.CreatedAt.UnixMilli(),
		BanUntil:   banUntil,
		BanReason:  banReason,
	}))
	// Create log
	c.Set("log", fmt.Sprintf("管理权限查询用户成功, 目标用户名: %s", req.UserName))
}
