package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/utils"
)

type AddressGeocodeByIDRequest struct {
	ID uint32 `uri:"id" binding:"required,gt=0"`
}

type AddressGeocodeByIDResponse struct {
	Result struct {
		Address *ShortGeocoderAddress
	} `json:"result"`
}

func (service *AddressService) AddressGeocodeByID(c *gin.Context) {
	request := AddressGeocodeByIDRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	result := service.geo.AddressByID(request.ID)
	if result == nil {
		utils.GErrorNotFound(c, "No such address")
		return
	}
	shortAddress, err := FullToShortAddress(result)
	if err != nil {
		utils.GErrorNotFound(c, err.Error())
		return
	}

	response := AddressGeocodeByIDResponse{}
	response.Result.Address = shortAddress

	utils.GSuccess(c, &response)
}
