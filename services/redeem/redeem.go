package redeem

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/services/slot"
	"bunker-web/services/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Generate(codeType int, note string, count int) []string {
	results := make([]string, count)
	for i := 0; i < count; {
		newCode := uuid.NewString()
		co := &models.RedeemCode{
			Code:     newCode,
			CodeType: codeType,
			Used:     false,
			Note:     note,
		}
		if models.DBCreate(co) != nil {
			continue
		}
		results[i] = newCode
		i++
	}
	return results
}

func UseRedeemCode(usr *models.User, code string) (string, *gin.Error) {
	// Get and check code status
	redeemCode, ginerr := getRedeemCode(code)
	if ginerr != nil {
		return "", ginerr
	}
	// Redeem
	var redeemResult string
	switch redeemCode.CodeType {
	case RedeemTypeUserOneMonth, RedeemTypeUserThreeMonth:
		// Get time
		time, ginerr := getRedeemTime(redeemCode)
		if ginerr != nil {
			return "", ginerr
		}
		// Set result
		redeemResult = generateHint(RedeemHintUser, time)
		// Active user
		if usr.Permission == user.PermissionGuest {
			usr.Permission = user.PermissionNormal
			redeemResult += " + 激活权限"
		}
		// Upgrade and renew user
		if ginerr := user.ExtendExpireTime(usr, time); ginerr != nil {
			return "", ginerr
		}
		// Update code status
		redeemCode.Used = true
		redeemCode.UserID = usr.ID
		// Save
		return redeemResult, giner.NewPrivateGinError(models.DBSave(redeemCode))
	case RedeemTypeSlotOneMonth, RedeemTypeSlotThreeMonth:
		// Get time
		time, ginerr := getRedeemTime(redeemCode)
		if ginerr != nil {
			return "", ginerr
		}
		// Create slot
		if _, ginerr := slot.Create(usr.ID, time); ginerr != nil {
			return "", ginerr
		}
		// Update code status
		redeemCode.Used = true
		redeemCode.UserID = usr.ID
		// Save
		return generateHint(RedeemHintSlot, time), giner.NewPrivateGinError(models.DBSave(redeemCode))
	default:
		return "", giner.NewPublicGinError("无效的兑换码类型")
	}
}

func UseRedeemCodeForSlot(usr *models.User, s *models.Slot, code string) (string, *gin.Error) {
	// Get and check code status
	redeemCode, ginerr := getRedeemCode(code)
	if ginerr != nil {
		return "", ginerr
	}
	if redeemCode.CodeType != RedeemTypeSlotOneMonth && redeemCode.CodeType != RedeemTypeSlotThreeMonth {
		return "", giner.NewPublicGinError("无法使用非 slot 类型的兑换码")
	}
	// Get time
	time, ginerr := getRedeemTime(redeemCode)
	if ginerr != nil {
		return "", ginerr
	}
	// Upgrade and renew slot
	if ginerr := slot.ExtendSlotExpireTime(s, time); ginerr != nil {
		return "", ginerr
	}
	// Update code status
	redeemCode.Used = true
	redeemCode.UserID = usr.ID
	return generateHint(RedeemHintSlot, time), giner.NewPrivateGinError(models.DBSave(redeemCode))
}
