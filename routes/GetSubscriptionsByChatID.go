package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type GetSubscriptionsByChatIDRequest struct {
	ChatID int64 `uri:"ChatID" binding:"required,gt=0"`
}

type GetSubscriptionsByChatIDResponse struct {
	Result []models.Subscription `json:"result"`
}

func (service *SubscriptionService) GetSubscriptionsByChatID(c *gin.Context) {
	var request GetSubscriptionsByChatIDRequest
	if err := c.ShouldBindUri(&request); err != nil {
		utils.GErrorBadRequest(c, err.Error())
		return
	}

	var subscriptions []models.Subscription
	service.db.
		Where(models.Subscription{ChatID: request.ChatID}).
		Find(&subscriptions)

	utils.GSuccess(c, &GetSubscriptionsByChatIDResponse{Result: subscriptions})
}
