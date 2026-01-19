package phoenix

import (
	"bunker-core/mcp"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransferCheckNumRequest struct {
	Data string `json:"data" binding:"min=1"`
}

type TransferCheckNumResponse struct {
	giner.BasicResponse
	Value string `json:"value"`
}

func (*Phoenix) TransferCheckNum(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	// Parse request
	var req TransferCheckNumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check if is phoenix
	if isPhoenix, ok := session.Load("isPhoenix"); !ok || !isPhoenix.(bool) {
		c.Error(giner.NewPublicGinError("无效会话"))
		return
	}
	// Get engineVersion
	engineVersion, ok := session.Load("engineVersion")
	if !ok {
		c.Error(giner.NewPublicGinError("会话已失效, 请重新登录"))
		return
	}
	// Get patchVersion
	patchVersion, ok := session.Load("patchVersion")
	if !ok {
		c.Error(giner.NewPublicGinError("会话已失效, 请重新登录"))
		return
	}
	// Get platform
	platform, ok := session.Load("platform")
	if !ok {
		c.Error(giner.NewPublicGinError("会话已失效, 请重新登录"))
		return
	}
	// Parse fb req
	var dataList []any
	if err := json.Unmarshal([]byte(req.Data), &dataList); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	if len(dataList) != 3 {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	mcpData, ok := dataList[0].(string)
	if !ok || mcpData == "" {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	salt, ok := dataList[1].(string)
	if !ok || salt == "" {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	uid, ok := dataList[2].(float64)
	if !ok || uid <= 0 {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Get check num
	result, err := mcp.GetMCPCheckNum(
		engineVersion.(string),
		patchVersion.(string),
		mcpData,
		salt,
		strconv.Itoa(int(uid)),
		platform.(string),
	)
	if err != nil {
		c.Error(giner.NewPrivateGinError(err))
		return
	}
	// Return result
	c.JSON(http.StatusOK, &TransferCheckNumResponse{
		BasicResponse: giner.BasicResponse{
			Success: true,
			Message: "ok",
		},
		Value: result,
	})
	// Create log
	c.Set("log", "CheckNum 获取成功")
}
