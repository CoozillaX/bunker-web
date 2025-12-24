package user

import (
	"bunker-web/configs"
	"bunker-web/models"
	"bunker-web/pkg/fbtoken"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/utils"
	"bunker-web/services/slot"
	"bunker-web/services/unlimited_rental_server"
	"bunker-web/services/user_ban_record"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	PermissionGuest  = iota // Can login to user center, but not allowed to use any functions
	PermissionNormal        // Allow login to rental server but with limitation (their own servers) and period limitation
	PermissionAdmin         // Allow manage other accounts
)

const (
	SMSCodeLimit = 5
)

func Create(username, password string, permission uint) (*models.User, *gin.Error) {
	if len(username) < 3 || len(username) > 12 {
		return nil, giner.NewPublicGinError("用户名长度不符合要求")
	}
	if !regexp.MustCompile("^[A-Za-z0-9]+$").MatchString(username) {
		return nil, giner.NewPublicGinError("无效用户名")
	}
	if password == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		return nil, giner.NewPublicGinError("密码不能为空")
	}
	if user, _ := QueryByUsername(username); user != nil {
		return nil, giner.NewPublicGinError("用户名已存在")
	}
	u := &models.User{
		Username:   username,
		Password:   utils.SHA256Hex([]byte(password + configs.USER_PSW_SALT)),
		Permission: permission,
		GameID:     0,
	}
	return u, giner.NewPrivateGinError(models.DBCreate(u))
}

func Remove(u *models.User) *gin.Error {
	// Remove helper
	if u.HelperMpayUser != nil {
		models.DBRemove(u.HelperMpayUser)
	}
	// Remove owner
	if u.OwnerMpayUser != nil {
		models.DBRemove(u.OwnerMpayUser)
	}
	// Remove slots
	slot.DeleteAllByUserID(u.ID)
	return giner.NewPrivateGinError(models.DBRemove(u))
}

func ExtendExpireTime(u *models.User, second int64) *gin.Error {
	if time.Now().After(u.ExpireAt.Time) {
		u.ExpireAt = sql.NullTime{
			Time:  time.Now().Add(time.Duration(second) * time.Second),
			Valid: true,
		}
	} else {
		u.ExpireAt = sql.NullTime{
			Time:  u.ExpireAt.Time.Add(time.Duration(second) * time.Second),
			Valid: true,
		}
	}
	return giner.NewPrivateGinError(models.DBSave(u))
}

func ExtendUnlimitedTime(u *models.User, second int64) *gin.Error {
	if time.Now().After(u.UnlimitedUntil.Time) {
		u.UnlimitedUntil = sql.NullTime{
			Time:  time.Now().Add(time.Duration(second) * time.Second),
			Valid: true,
		}
	} else {
		u.UnlimitedUntil = sql.NullTime{
			Time:  u.UnlimitedUntil.Time.Add(time.Duration(second) * time.Second),
			Valid: true,
		}
	}
	return giner.NewPrivateGinError(models.DBSave(u))
}

func CheckIfVaild(u *models.User) (bool, string) {
	// Check ban
	if banRecord, _ := user_ban_record.GetCurrentBanRecordFormattedStringByUserID(u.ID); len(banRecord) > 0 {
		return false, banRecord
	}
	// Check permission
	switch u.Permission {
	case PermissionGuest:
		// EVENT
		// return false, "账户未激活"
	case PermissionNormal:
		// DISABLED
		// // Check if unlimited and expire time
		// if time.Now().After(u.ExpireAt.Time) && time.Now().After(u.UnlimitedUntil.Time) {
		// 	return false, "账户不在有效期内"
		// }
	}
	return true, ""
}

func QueryByUsername(username string) (*models.User, *gin.Error) {
	var user models.User
	err := models.DB.
		Preload("HelperMpayUser").
		Preload("OwnerMpayUser").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return &user, nil
}

func QueryByToken(token, currentHashedIP string) (*models.User, *gin.Error) {
	// 1. Parse token
	username, saltedPassword, hashedIP, err := fbtoken.Decrypt(token)
	if err != nil {
		return nil, giner.NewPublicGinError("无效的Token")
	}
	// 2. Check password
	usr, ginerr := NormalLogin(username, saltedPassword)
	if ginerr != nil {
		return nil, giner.NewPublicGinError("Token已失效, 请重新获取")
	}
	// 3. Check IP
	if hashedIP != "" && currentHashedIP != hashedIP {
		return nil, giner.NewPublicGinError("当前IP无法使用此Token")
	}
	return usr, nil
}

func QueryUserByAPIKey(key string) (*models.User, *gin.Error) {
	var user models.User
	err := models.DB.
		Preload("HelperMpayUser").
		Preload("OwnerMpayUser").
		Where("api_key = ?", key).
		First(&user).Error
	if err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return &user, nil
}

func QueryUserByEmail(email string) (*models.User, *gin.Error) {
	var user models.User
	err := models.DB.
		Preload("HelperMpayUser").
		Preload("OwnerMpayUser").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return &user, nil
}

func QueryUserByID(id uint) (*models.User, *gin.Error) {
	var user models.User
	err := models.DB.
		Preload("HelperMpayUser").
		Preload("OwnerMpayUser").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return &user, nil
}

func NormalLogin(username, saltedPassword string) (*models.User, *gin.Error) {
	usr, _ := QueryByUsername(username)
	if usr == nil || saltedPassword != usr.Password {
		return nil, giner.SetTranslationCode(giner.NewPublicGinError("无效的用户中心用户名或密码"), giner.C_Auth_InvalidUser)
	}
	return usr, nil
}

func PhoenixLogin(ip, token, username, hashedPassword string) (usr *models.User, ginerr *gin.Error) {
	// Token or password login
	if token != "" {
		if usr, ginerr = QueryByToken(token, utils.MD5Hex([]byte(ip))); ginerr != nil {
			return nil, giner.SetTranslationCode(ginerr, giner.C_Auth_InvalidToken)
		}
	} else {
		if usr, ginerr = NormalLogin(username, utils.SHA256Hex([]byte(hashedPassword+configs.USER_PSW_SALT))); ginerr != nil {
			return nil, ginerr
		}
	}
	// Check if vaild
	if vaild, reason := CheckIfVaild(usr); !vaild {
		return usr, giner.NewPublicGinError(reason)
	}
	return usr, nil
}

func GameLicenseCheck(u *models.User, serverCode string, ownerID int) *gin.Error {
	// Check user if admin
	if u.Permission == PermissionAdmin {
		return nil
	}
	// Check user if unlimited
	if u.Permission < PermissionAdmin && time.Now().Before(u.UnlimitedUntil.Time) {
		return nil
	}
	// Check server if unlimited
	if _, err := unlimited_rental_server.QueryByServerCode(serverCode); err == nil {
		return nil
	}
	// Check if not bind owner
	if u.GameID == 0 {
		u.GameID = ownerID
		return giner.NewPrivateGinError(models.DBSave(u))
	}
	// Check if owner
	if u.GameID == ownerID {
		return nil
	}
	// Check if slot
	if err := slot.CheckIfVaild(u.ID, ownerID); err != nil {
		return giner.NewPublicGinError(
			fmt.Sprintf(
				"登录失败, %v",
				err.Error(),
			),
		)
	}
	return nil
}
