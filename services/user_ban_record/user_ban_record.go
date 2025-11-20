package user_ban_record

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Create(userID uint, reason string, second int64) *gin.Error {
	banRecord := &models.UserBanRecord{
		UserID: userID,
		Until:  sql.NullTime{Time: time.Now().Add(time.Second * time.Duration(second)), Valid: true},
		Reason: reason,
	}
	return giner.NewPrivateGinError(models.DBCreate(banRecord))
}

func GetCurrentBanRecordByUserID(userID uint) (*models.UserBanRecord, *gin.Error) {
	var banRecord models.UserBanRecord
	result := models.DB.Where("user_id = ? AND until > ?", userID, time.Now()).First(&banRecord)
	if result.Error != nil {
		return nil, giner.NewPrivateGinError(result.Error)
	}
	return &banRecord, nil
}

func GetCurrentBanRecordFormattedStringByUserID(userID uint) (string, *gin.Error) {
	banRecord, ginerr := GetCurrentBanRecordByUserID(userID)
	if ginerr != nil {
		return "", ginerr
	}
	return fmt.Sprintf(
		"账户因 %s 被封禁至 %s",
		banRecord.Reason,
		banRecord.Until.Time.Format(time.DateTime),
	), nil
}

func RemoveAllBanRecordByUserID(userID uint) *gin.Error {
	result := models.DB.Where("user_id = ?", userID).Delete(&models.UserBanRecord{})
	return giner.NewPrivateGinError(result.Error)
}
