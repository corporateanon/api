package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type AddressGetTotalCountResponse struct {
	Result int64 `json:"result"`
}

func (service *AddressService) AddressGetTotalCount(c *gin.Context) {
	var cnt int64 = 0
	err := service.db.Model(&models.AddressAr{}).Count(&cnt).Error
	if err != nil {
		utils.GErrorInternal(c, err.Error())
		return
	}
	utils.GSuccess(c, &AddressGetTotalCountResponse{Result: cnt})
}
