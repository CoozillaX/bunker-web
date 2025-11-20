package notice

import (
	"bunker-web/pkg/giner"
	"bunker-web/services/announcement"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	PageNum  int `json:"page_num" binding:"gt=0"`
	PageSize int `json:"page_size" binding:"gt=0"`
}

type NoticeData struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorName string `json:"author_name"`
	CreateAt   int64  `json:"create_at"`
}

type QueryResponseData struct {
	Total   int64         `json:"total"`
	Notices []*NoticeData `json:"notices"`
}

func (*Notice) Query(c *gin.Context) {
	// Parse request
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Limit page size
	if req.PageSize > 10 {
		req.PageSize = 10
	}
	// Query notice count
	total, ginerr := announcement.GetTotal()
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Query notice by page
	result, ginerr := announcement.QueryByPage(req.PageNum, req.PageSize)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Format response
	var noticeData []*NoticeData
	for _, item := range result {
		noticeData = append(noticeData, &NoticeData{
			ID:         item.ID,
			Title:      item.Title,
			Content:    item.Content,
			AuthorName: item.AuthorName,
			CreateAt:   item.CreatedAt.UnixMilli(),
		})
	}
	c.JSON(http.StatusOK, giner.MakeHTTPResponse(true).SetData(&QueryResponseData{
		Total:   total,
		Notices: noticeData,
	}))
	// Create log
	c.Set("log", "公告查询成功")
}
