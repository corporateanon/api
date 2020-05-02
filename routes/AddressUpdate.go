package routes

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type AddressUpdateRequestBody struct {
	CheckStatus    models.AddressArCheckStatus `json:"CheckStatus" binding:"required"`
	ServiceMessage string                      `json:"ServiceMessage" binding:"required"`
	Hash           string                      `json:"Hash" binding:"required"`
}
type AddressUpdateRequestURI struct {
	ID int64 `uri:"id" binding:"required,gt=0"`
}
type AddressUpdateResponse struct {
	Result *models.AddressAr `json:"result"`
}

func (service *AddressService) AddressUpdate(c *gin.Context) {
	request := AddressUpdateRequestBody{}
	if err := c.ShouldBind(&request); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}
	uriParams := AddressUpdateRequestURI{}
	if err := c.ShouldBindUri(&uriParams); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	if request.CheckStatus != models.AddressStatusNoWork && request.CheckStatus != models.AddressStatusWork {
		utils.GErrorBadRequest(c, fmt.Sprintf("CheckStatus may be either '%s' or '%s'", models.AddressStatusNoWork, models.AddressStatusWork))
		return
	}

	address := &models.AddressAr{}

	if err := service.db.
		Where(uriParams.ID).
		Preload("Subscriptions").
		First(address).Error; err != nil {
		utils.GErrorInternal(c, err.Error())
		return
	}
	if address.ID == 0 {
		utils.GErrorNotFound(c, "Not found")
		return
	}

	previousHash := address.Hash

	address.CheckStatus = request.CheckStatus
	address.ServiceMessage = request.ServiceMessage
	address.Hash = request.Hash
	address.CheckedAt = time.Now()
	if err := service.db.Save(address).Error; err != nil {
		utils.GErrorInternal(c, err.Error())
		return
	}
	utils.GSuccess(c, &AddressUpdateResponse{Result: address})

	//------------

	needsNotification := previousHash != request.Hash &&
		address.Subscriptions != nil &&
		len(address.Subscriptions) > 0

	if needsNotification {
		chatIDs := make([]int64, len(address.Subscriptions))
		for i, subscription := range address.Subscriptions {
			chatIDs[i] = subscription.ChatID
		}
		geocoderAddress := service.geo.AddressByID(uint32(address.ID))
		service.notifier.NotifyServiceMessageChange(
			chatIDs,
			request.ServiceMessage,
			formatAddress(geocoderAddress),
			address.CheckStatus,
		)
	}
}
