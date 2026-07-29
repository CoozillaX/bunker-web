package other

import (
	"bunker-web/models"
	"bunker-web/pkg/g79"
	"bunker-web/pkg/giner"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RandomNameResponse struct {
	giner.BasicResponse
	Data []g79.RandomName `json:"data"` // 随机昵称候选列表
}

// RandomName godoc
//
//	@Summary		获取随机昵称
//	@Description	获取游戏生成的随机昵称候选列表, 需要提供 API Key 并创建辅助用户
//	@Description	注意: 此 API 会尝试进行游戏登录 (Helper)
//	@Tags			其他
//	@Produce		json
//	@Param			Authorization	header		string				true	"此处需要填写 API Key, 可以从用户菜单获取"
//	@Success		200				{object}	RandomNameResponse	"成功时返回"
//	@Router			/openapi/other/random_name [post]
func (*Other) RandomName(c *gin.Context) {
	u, _ := c.Get("usr")
	usr := u.(*models.User)
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("未创建辅助用户"))
		return
	}

	gu, ginerr := g79.HandleG79Login(usr.HelperMpayUser)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	names, ginerr := g79.GetRandomNames(gu)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}

	c.JSON(http.StatusOK, &RandomNameResponse{
		BasicResponse: giner.BasicResponse{Success: true, Message: "ok"},
		Data:          names,
	})
	c.Set("log", "获取随机昵称成功")
}
