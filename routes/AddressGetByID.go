package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type AddressGetByIDRequest struct {
	ID int64 `uri:"id" binding:"required,gt=0"`
}
type AddressGetByIDResponse struct {
	Result *models.AddressAr `json:"result"`
}

func (service *AddressService) AddressGetByID(c *gin.Context) {
	request := AddressGetByIDRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	address := &models.AddressAr{}

	service.db.
		Preload("Subscriptions").
		Where(request.ID).
		First(address)
	if address.ID == 0 {
		utils.GErrorNotFound(c, "Address not found")
		return
	}

	response := AddressGetByIDResponse{
		Result: address,
	}

	utils.GSuccess(c, response)
}
