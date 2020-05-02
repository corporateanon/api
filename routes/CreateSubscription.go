package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type CreateSubscriptionRequest struct {
	ChatID      int64 `json:"ChatID"      binding:"required,gt=0"`
	AddressArID int64 `json:"AddressArID" binding:"required,gt=0"`
}

type CreateSubscriptionResponse struct {
	Result *models.Subscription `json:"result"`
}

func (service *SubscriptionService) CreateSubscription(c *gin.Context) {
	var payload CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	addressFromGeocoder := service.geo.AddressByID(uint32(payload.AddressArID))
	if addressFromGeocoder == nil {
		utils.GErrorBadRequest(c, "AddressArID not found")
		return
	}

	subscription := &models.Subscription{
		ChatID:      payload.ChatID,
		AddressArID: payload.AddressArID,
	}

	existingSubscription := &models.Subscription{}

	tx := service.db.Begin()

	existingAddress := &models.AddressAr{}
	tx.Where(subscription.AddressArID).First(existingAddress)
	if existingAddress.ID == 0 {
		existingAddress.CheckStatus = models.AddressStatusInit
		existingAddress.ID = uint(subscription.AddressArID)
		tx.Save(existingAddress)
	}

	tx.Where(subscription).First(existingSubscription)
	if existingSubscription.ID > 0 {
		utils.GSuccess(c, &CreateSubscriptionResponse{Result: existingSubscription})
		tx.Commit()
		return
	}

	tx.Save(subscription)

	if err := tx.Commit().Error; err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	utils.GSuccess(c, &CreateSubscriptionResponse{Result: subscription})
}
