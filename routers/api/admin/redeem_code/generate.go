package redeem_code

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/redeem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GenerateRequest struct {
	CodeType int    `json:"type" binding:"gte=1,lte=4"`
	Note     string `json:"note" binding:"min=1"`
	Count    int    `json:"count" binding:"gte=1,lte=999"`
}

func (*RedeemCode) Generate(c *gin.Context) {
	// Parse request
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Generate redeem code
	codes := redeem.Generate(req.CodeType, req.Note, req.Count)
	// Setup file
	filename := fmt.Sprintf("redeem_%d_%d_%s.txt", req.CodeType, req.Count, time.Now().Format("2006-01-02_15-04-05"))
	fileContent := strings.Join(codes, "\n")
	// Return as file
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.String(http.StatusOK, fileContent)
	// Create log
	c.Set("log", fmt.Sprintf("管理权限生成兑换码成功, 类型: %d, 数量: %d", req.CodeType, req.Count))
}
