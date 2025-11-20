package unlimited_rental_server

import (
	"bunker-web/models"
	"bunker-web/pkg/giner"

	"github.com/gin-gonic/gin"
)

func Create(userID uint, serverCode string) (*models.UnlimitedRentalServer, *gin.Error) {
	s := &models.UnlimitedRentalServer{
		OperatorID: userID,
		ServerCode: serverCode,
	}
	return s, giner.NewPrivateGinError(models.DBCreate(s))
}

func Delete(s *models.UnlimitedRentalServer) *gin.Error {
	return giner.NewPrivateGinError(models.DBDelete(s))
}

func QueryByID(id uint) (*models.UnlimitedRentalServer, *gin.Error) {
	s := &models.UnlimitedRentalServer{}
	err := models.DB.Where("id = ?", id).First(s).Error
	return s, giner.NewPrivateGinError(err)
}

func QueryByServerCode(serverCode string) (*models.UnlimitedRentalServer, *gin.Error) {
	s := &models.UnlimitedRentalServer{}
	err := models.DB.Where("server_code = ?", serverCode).First(s).Error
	return s, giner.NewPrivateGinError(err)
}

func QueryAll() ([]models.UnlimitedRentalServer, *gin.Error) {
	s := []models.UnlimitedRentalServer{}
	err := models.DB.Find(&s).Error
	return s, giner.NewPrivateGinError(err)
}
