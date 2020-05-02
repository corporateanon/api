package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/utils"
)

type AddressGeocodeRequest struct {
	Lat      float64 `uri:"lat"`
	Lng      float64 `uri:"lng"`
	Accuracy float64 `uri:"accuracy"`
}

type AddressGeocodeResponse struct {
	Result GeocodeResult `json:"result"`
}

type GeocodeResponseAddress struct {
	ID            uint32
	Distance      float64
	AddressString string
}
type GeocodeResult struct {
	Addresses []GeocodeResponseAddress
}

func (service *AddressService) AddressGeocode(c *gin.Context) {
	request := AddressGeocodeRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	result := service.geo.ReverseGeocode(request.Lat, request.Lng, request.Accuracy, 10)

	response := AddressGeocodeResponse{}
	for _, item := range result {
		response.Result.Addresses = append(response.Result.Addresses, GeocodeResponseAddress{
			ID:            item.FullAddress.Address.ID,
			Distance:      item.Distance,
			AddressString: formatGeocodingResult(item),
		})
	}
	utils.GSuccess(c, &response)
}
