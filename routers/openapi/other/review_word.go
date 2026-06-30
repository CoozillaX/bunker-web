package other

import (
	"bunker-web/pkg/giner"
	"context"
	"net/http"
	"time"

	"github.com/CoozillaX/g79-regex/review"
	"github.com/gin-gonic/gin"
)

type ReviewResult struct {
	Group string `json:"group"` // 匹配的敏感词规则组, 其中 "nickname" 只在昵称检测时使用
	Index int    `json:"index"` // 匹配的敏感词在文本中的起始位置
	Start int    `json:"start"` // 匹配的敏感词在文本中的起始位置
	End   int    `json:"end"`   // 匹配的敏感词在文本中的结束位置
	Text  string `json:"text"`  // 匹配的敏感词
}

type ReviewWordRequest struct {
	Content   string `json:"content" binding:"min=1"`                          // 被检测的文本
	Level     string `json:"level" binding:"max=100" example:"0"`              // 检测级别, 不填写则使用默认值 0
	Channel   string `json:"channel" binding:"max=100" example:"item_comment"` // 检测频道, 不填写则使用默认值 item_comment
	FirstOnly bool   `json:"first_only" example:"false"`                       // 是否只返回第一个匹配的敏感词, 默认为 false
}

type ReviewWordResponse struct {
	giner.BasicResponse
	Data []ReviewResult `json:"data"`
}

var (
	reviewer       *review.Reviewer
	lastUpdateTime time.Time
)

// ReviewWord godoc
//
//	@Summary		敏感词检查
//	@Description	检查文本中是否包含敏感词, 无需提供API key, 敏感词库每小时更新一次<br>
//	@Description	以下为额外的字段说明: <br>
//	@Description	level: 默认使用 0, 少部分表达式会使用此字段<br>
//	@Description	channel: 在不同的场景下此值会不一样, 部分值及其说明如下<br>
//	@Description	- item_comment: 默认情况下使用此channel<br>
//	@Description	- check_long_numbers: 需要额外检查文本是否包含长数字时使用此channel<br>
//	@Description	- World: 文本来源为世界聊天时使用此channel<br>
//	@Tags			其他
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		ReviewWordRequest		true	"请求时需要在 Body 携带以下查询参数"
//	@Success		200		{object}	ReviewWordResponse	"成功时返回"
//	@Router			/openapi/other/review_text [post]
func (*Other) ReviewWord(c *gin.Context) {
	// Parse request
	var req ReviewWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Check if reviewer need to be initialised or updated
	if reviewer == nil || time.Since(lastUpdateTime) > time.Hour {
		var err error

		if reviewer == nil {
			reviewer, err = review.New(context.Background())
		} else {
			_, err = reviewer.Reload(context.Background())
		}

		if err != nil {
			c.Error(giner.NewPublicGinError("初始化/更新敏感词检查器时出现问题"))
			return
		}
		lastUpdateTime = time.Now()
	}
	// Review content
	result := reviewer.ReviewWord(req.Content, &review.Options{
		Level:     req.Level,
		Channel:   req.Channel,
		FirstOnly: req.FirstOnly,
	})
	var data []ReviewResult
	for _, r := range result {
		data = append(data, ReviewResult{
			Group: r.Group,
			Index: r.Index,
			Start: r.Start,
			End:   r.End,
			Text:  r.Text,
		})
	}
	// Pack player info
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(data))
	// Create log
	c.Set("log", "敏感词检查成功")
}
