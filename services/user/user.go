package user

import (
	"bunker-web/configs"
	"bunker-web/models"
	"bunker-web/pkg/fbtoken"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/utils"
	"bunker-web/services/unlimited_rental_server"
	"bunker-web/services/user_ban_record"
	"regexp"

	"github.com/gin-gonic/gin"
)

const (
	PermissionNormal    = iota // Allow login to rental server but with limitation (their own servers) and period limitation
	PermissionDeveloper        // Allow login to rental server without limitation
	PermissionAdmin            // Allow manage other accounts
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
	return giner.NewPrivateGinError(models.DBRemove(u))
}

func CheckIfVaild(u *models.User) (bool, string) {
	// Check ban
	if banRecord, _ := user_ban_record.GetCurrentBanRecordFormattedStringByUserID(u.ID); len(banRecord) > 0 {
		return false, banRecord
	}
	return true, ""
}

func QueryByUsername(username string) (*models.User, *gin.Error) {
	var user models.User
	err := models.DB.
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return &user, nil
}

func QueryByToken(token, currentHashedIP string) (*models.User, *gin.Error) {
	// 1. Parse token
	username, saltedPassword, err := fbtoken.Decrypt(token)
	if err != nil {
		return nil, giner.NewPublicGinError("无效的Token")
	}
	// 2. Check password
	usr, ginerr := NormalLogin(username, saltedPassword)
	if ginerr != nil {
		return nil, giner.NewPublicGinError("Token已失效, 请重新获取")
	}
	return usr, nil
}

func QueryUserByAPIKey(key string) (*models.User, *gin.Error) {
	var user models.User
	err := models.DB.
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
	// Check user if developer or admin
	if u.Permission >= PermissionDeveloper {
		return nil
	}
	// Check server if unlimited
	if _, err := unlimited_rental_server.QueryByServerCode(serverCode); err == nil {
		return nil
	}
	// Check if not bind owner
	if u.GameID == 0 {
		return giner.NewPublicGinError("请前往用户中心绑定游戏ID后再尝试登录")
	}
	// Check if owner
	if u.GameID == ownerID {
		return nil
	}
	return giner.NewPublicGinError("登录失败，此服务器拥有者的游戏ID与当前绑定的游戏ID不匹配")
}

func DeleteHelper(u *models.User) *gin.Error {
	if u.HelperMpayUser == nil {
		return nil
	}
	// Remove helper
	if err := models.DBRemove(u.HelperMpayUser); err != nil {
		return giner.NewPrivateGinError(err)
	}
	// Reset helper info
	u.HelperMpayUser = nil
	return giner.NewPrivateGinError(models.DBSave(u))
}

func DeleteOwner(u *models.User) *gin.Error {
	if u.OwnerMpayUser == nil {
		return nil
	}
	// Remove helper
	if err := models.DBRemove(u.OwnerMpayUser); err != nil {
		return giner.NewPrivateGinError(err)
	}
	// Reset helper info
	u.OwnerMpayUser = nil
	return giner.NewPrivateGinError(models.DBSave(u))
}
