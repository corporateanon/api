package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type AddressGetListRequest struct {
}

type ExtendedAddress struct {
	models.AddressAr
	AddressDetails *ShortGeocoderAddress
}

type AddressGetListResponse struct {
	Result []ExtendedAddress `json:"result"`
}

func (service *AddressService) AddressGetList(c *gin.Context) {
	listQueryParams, err := utils.DecodeReactAdminQueryParams(c)
	if err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	addresses := []models.AddressAr{}

	if err := models.ApplyListQueryParams(
		service.db,
		models.AddressAr{},
		listQueryParams,
	).Find(&addresses).Error; err != nil {
		utils.GErrorInternal(c, err.Error())
		return
	}

	if contentRange, err := models.GetResultRange(
		service.db,
		&models.AddressAr{},
		listQueryParams,
		"address",
	); err != nil {
		utils.GErrorInternal(c, err.Error())
		return
	} else if contentRange != "" {
		c.Header("Content-Range", contentRange)
	}

	extendedAddresses := make([]ExtendedAddress, len(addresses))
	for i, address := range addresses {

		fullAddress := service.geo.AddressByID(uint32(address.ID))

		if fullAddress == nil {
			extendedAddresses[i] = ExtendedAddress{
				AddressAr: address,
			}
			continue
		}

		shortAddress, err := FullToShortAddress(fullAddress)
		if err != nil {
			extendedAddresses[i] = ExtendedAddress{
				AddressAr: address,
			}
			continue
		}
		extendedAddresses[i] = ExtendedAddress{
			AddressAr:      address,
			AddressDetails: shortAddress,
		}
	}

	utils.GSuccess(c, &AddressGetListResponse{Result: extendedAddresses})
}
