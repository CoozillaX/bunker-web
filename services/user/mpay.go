package user

import (
	"bunker-core/protocol/defines"
	"bunker-core/protocol/mpay"
	"bunker-web/models"
	"bunker-web/pkg/giner"

	"github.com/gin-gonic/gin"
)

// This function will ensure the helper mpay user is created and persisted.
func GetLoginHelperForHelper(u *models.User) (*mpay.LoginHelper, *gin.Error) {
	// Initialise Mpay user helper
	var mu *defines.MpayUser
	if u.HelperMpayUser != nil && u.HelperMpayUser.MpayUser != nil {
		mu = u.HelperMpayUser.MpayUser
	}
	helper := mpay.CreateLoginHelper(mu)
	// Initialise helper mpay user if not exist
	if u.HelperMpayUser == nil {
		u.HelperMpayUser = &models.MpayUser{
			MpayUser: helper.GetMpayUser(),
		}
		// Store to DB
		if err := models.DBCreate(u.HelperMpayUser); err != nil {
			return nil, giner.NewPrivateGinError(err)
		}
		if err := models.DBSave(u); err != nil {
			return nil, giner.NewPrivateGinError(err)
		}
	}
	return helper, nil
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
	u.HelperMpayUserID = nil
	u.HelperMpayUser = nil
	return giner.NewPrivateGinError(models.DBSave(u))
}

// This function will ensure the helper mpay user is created and persisted.
func GetLoginHelperForOwner(u *models.User) (*mpay.LoginHelper, *gin.Error) {
	// Initialise Mpay user helper
	var mu *defines.MpayUser
	if u.OwnerMpayUser != nil && u.OwnerMpayUser.MpayUser != nil {
		mu = u.OwnerMpayUser.MpayUser
	}
	helper := mpay.CreateLoginHelper(mu)
	// Initialise helper mpay user if not exist
	if u.OwnerMpayUser == nil {
		u.OwnerMpayUser = &models.MpayUser{
			MpayUser: helper.GetMpayUser(),
		}
		// Store to DB
		if err := models.DBCreate(u.OwnerMpayUser); err != nil {
			return nil, giner.NewPrivateGinError(err)
		}
		if err := models.DBSave(u); err != nil {
			return nil, giner.NewPrivateGinError(err)
		}
	}
	return helper, nil
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
	u.OwnerMpayUserID = nil
	u.OwnerMpayUser = nil
	return giner.NewPrivateGinError(models.DBSave(u))
}
