package redeem

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"fmt"

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

func generateHint(hint string, duration int64) string {
	return fmt.Sprintf(hint, duration/(60*60*24))
}

func getRedeemCode(code string) (*models.RedeemCode, *gin.Error) {
	var redeemCode models.RedeemCode
	if query := models.DB.Where("code = ?", code).First(&redeemCode); query.Error != nil {
		return nil, giner.NewPublicGinError("兑换码无效")
	}
	if redeemCode.Used {
		return nil, giner.NewPublicGinError("兑换码已被使用")
	}
	return &redeemCode, nil
}

func getRedeemTime(redeem *models.RedeemCode) (int64, *gin.Error) {
	time, found := RedeemTimeMap[redeem.CodeType]
	if !found {
		return 0, giner.NewPublicGinError("无效的兑换码类型")
	}
	return time, nil
}
