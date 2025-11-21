package redeem

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RedeemTypeUserOneMonth = iota + 1
	RedeemTypeUserThreeMonth
	RedeemTypeSlotOneMonth
	RedeemTypeSlotThreeMonth
)

const (
	RedeemHintUser = "兑换成功, 您已获得 %d 天账户有效期"
	RedeemHintSlot = "兑换成功, 您已获得 %d 天 slot 有效期"
)

const (
	RedeemTimeOneMonth   = 60 * 60 * 24 * 31 * 1
	RedeemTimeThreeMonth = 60 * 60 * 24 * 31 * 3
	RedeemTimeOneYear    = 60 * 60 * 24 * 31 * 12
	RedeemTimeTenYear    = 60 * 60 * 24 * 31 * 12 * 10
)

var RedeemTimeMap = map[int]int64{
	RedeemTypeUserOneMonth:   RedeemTimeOneMonth,
	RedeemTypeUserThreeMonth: RedeemTimeThreeMonth,
	RedeemTypeSlotOneMonth:   RedeemTimeOneMonth,
	RedeemTypeSlotThreeMonth: RedeemTimeThreeMonth,
}

func isEventTime() bool {
	// Timezone: Asia/Shanghai (UTC+8)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	end := time.Date(2024, 11, 12, 23, 59, 59, 0, loc)
	// Check if now is in event time
	return time.Now().Before(end)
}

func generateHint(hint string, duration int64) string {
	return fmt.Sprintf(hint, duration/(60*60*24))
}

func getRedeemCode(_ string) (*models.RedeemCode, *gin.Error) {
	return nil, giner.NewPublicGinError("兑换码功能已被禁用")
	// var redeemCode models.RedeemCode
	// if query := models.DB.Where("code = ?", code).First(&redeemCode); query.Error != nil {
	// 	return nil, giner.NewPublicGinError("兑换码无效")
	// }
	// if redeemCode.Used {
	// 	return nil, giner.NewPublicGinError("兑换码已被使用")
	// }
	// return &redeemCode, nil
}

func getRedeemTime(redeem *models.RedeemCode) (int64, *gin.Error) {
	time, found := RedeemTimeMap[redeem.CodeType]
	if !found {
		return 0, giner.NewPublicGinError("无效的兑换码类型")
	}
	// EVENT
	if isEventTime() {
		// 130%
		time = int64(float64(time) * 1.3)
	}
	return time, nil
}
