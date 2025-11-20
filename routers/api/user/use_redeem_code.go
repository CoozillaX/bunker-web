package user

import (
	"bunker-web/pkg/giner"

	"github.com/gin-gonic/gin"
)

type RedeemRequest struct {
	Code string `json:"redeem_code" binding:"len=36"`
}

func (*User) Redeem(c *gin.Context) {
	// DISABLED
	c.Error(giner.NewPublicGinError("兑换码功能已被禁用"))

	// // Get session
	// bearer, _ := c.Get("bearer")
	// session, _ := sessions.GetSessionByBearer(bearer.(string))
	// u, _ := session.Load("usr")
	// usr := u.(*models.User)
	// // Parse request
	// var req RedeemRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.Error(giner.NewPublicGinError("无效参数"))
	// 	return
	// }
	// // Use redeem code
	// redeemResult, ginerr := redeem.UseRedeemCode(usr, req.Code)
	// if ginerr != nil {
	// 	c.Error(ginerr)
	// 	return
	// }
	// c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetMessage(redeemResult))
	// // Create log
	// c.Set("log", redeemResult)
}
