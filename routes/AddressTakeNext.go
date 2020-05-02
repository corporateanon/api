package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type AddressTakeNextResponse struct {
	Result struct {
		Address         *models.AddressAr
		GeocoderAddress *ShortGeocoderAddress
	} `json:"result"`
}

func (service *AddressService) AddressTakeNext(c *gin.Context) {
	address := &models.AddressAr{}
	err := service.db.Order("taken_at ASC").First(address).Error
	if err != nil {
		if err.Error() == "record not found" {
			utils.GErrorNotFound(c, "No addresses in database")
			return
		}
		utils.GErrorInternal(c, err.Error())
	}
	if address.ID == 0 {
		utils.GErrorNotFound(c, "No addresses in database")
		return
	}

	address.TakenAt = time.Now()
	err = service.db.Save(address).Error
	if err != nil {
		utils.GErrorInternal(c, err.Error())
		return
	}

	geocoderAddress := service.geo.AddressByID(uint32(address.ID))
	if geocoderAddress == nil {
		utils.GErrorNotFound(c, "Address does not exist in geocoder")
		return
	}

	short, err := FullToShortAddress(geocoderAddress)
	if err != nil {
		utils.GErrorNotFound(c, err.Error())
		return
	}

	response := AddressTakeNextResponse{}
	response.Result.Address = address
	response.Result.GeocoderAddress = short
	utils.GSuccess(c, &response)
}
