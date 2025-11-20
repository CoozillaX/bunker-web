package slot

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Create(userID uint, second int64) (*models.Slot, *gin.Error) {
	s := &models.Slot{
		UserID: userID,
		ExpireAt: sql.NullTime{
			Time:  time.Now().Add(time.Duration(second) * time.Second),
			Valid: true,
		},
	}
	return s, giner.NewPrivateGinError(models.DBCreate(s))
}

func Delete(s *models.Slot) *gin.Error {
	return giner.NewPrivateGinError(models.DBRemove(s))
}

func DeleteAllByUserID(userID uint) *gin.Error {
	return giner.NewPrivateGinError(models.DB.Where("user_id = ?", userID).Delete(&models.Slot{}).Error)
}

func SetGameID(s *models.Slot, userID uint, gameID int, note string) *gin.Error {
	if s.GameID != 0 {
		return giner.NewPublicGinError("不允许重复设置游戏ID")
	}
	query := &models.Slot{}
	if err := models.DB.Where("user_id = ? AND game_id = ?", userID, gameID).First(query).Error; err == nil {
		return giner.NewPublicGinError(
			fmt.Sprintf("已存在绑定当前服务器拥有者ID (%d) 的 slot, 无需重复绑定", query.GameID),
		)
	}
	s.GameID = gameID
	s.Note = note
	s.ExpireAt = sql.NullTime{
		Time: time.Now().Add(
			s.ExpireAt.Time.Sub(s.CreatedAt),
		),
		Valid: true,
	}
	return giner.NewPrivateGinError(models.DBSave(s))
}

func ExtendSlotExpireTime(s *models.Slot, second int64) *gin.Error {
	if time.Now().After(s.ExpireAt.Time) {
		s.ExpireAt = sql.NullTime{
			Time:  time.Now().Add(time.Duration(second) * time.Second),
			Valid: true,
		}
	} else {
		s.ExpireAt = sql.NullTime{
			Time:  s.ExpireAt.Time.Add(time.Duration(second) * time.Second),
			Valid: true,
		}
	}
	return giner.NewPrivateGinError(models.DBSave(s))
}

func QueryByID(id uint) (*models.Slot, *gin.Error) {
	s := &models.Slot{}
	err := models.DB.Where("id = ?", id).First(s).Error
	return s, giner.NewPrivateGinError(err)
}

func QuerySlotListByUserID(userID uint) (slots []models.Slot) {
	models.DB.Where("user_id = ?", userID).Find(&slots)
	return slots
}

func CheckIfVaild(userID uint, gameID int) error {
	s := &models.Slot{}
	err := models.DB.Where("user_id = ? AND game_id = ?", userID, gameID).First(s).Error
	if err != nil {
		return fmt.Errorf("需要有效的 slot 来进入非绑定游戏ID的服务器")
	}
	if !time.Now().Before(s.ExpireAt.Time) {
		return fmt.Errorf("对应的 slot 不在有效期内")
	}
	return nil
}
