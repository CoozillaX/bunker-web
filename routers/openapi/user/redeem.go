package user

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/redeem"
	"bunker-web/services/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedeemRequest struct {
	Username string `json:"username" binding:"min=1"`     // 用户名
	Code     string `json:"redeem_code" binding:"len=36"` // 要使用的兑换码
}

type RedeemResponseData struct {
	Result string `json:"result"` // 兑换结果
}

type RedeemResponse struct {
	giner.BasicResponse
	Data *RedeemResponseData `json:"data"`
}

// Redeem godoc
//
//	@Summary		使用兑换码
//	@Description	使用兑换码, 无需进行登录验证
//	@Tags			用户中心
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		RedeemRequest	true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200		{object}	RedeemResponse	"成功时返回"
//	@Router			/openapi/user/redeem [post]
func (*User) Redeem(c *gin.Context) {
	// Parse request
	var req RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get user
	usr, ginerr := user.QueryByUsername(req.Username)
	if ginerr != nil {
		c.Error(giner.NewPublicGinError("兑换失败"))
		return
	}
	// Use redeem code
	redeemResult, ginerr := redeem.UseRedeemCode(usr, req.Code)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(&RedeemResponseData{
		Result: redeemResult,
	}))
	// Create log
	c.Set("log", fmt.Sprintf("API兑换成功, 用户名 %s: %s", req.Username, redeemResult))
}
