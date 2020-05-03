package router

import (
	"github.com/gin-gonic/gin"
	"github.com/my1562/api/routes"
)

func NewGinRouter(
	subscriptionService *routes.SubscriptionService,
	addressService *routes.AddressService,
) *gin.Engine {
	route := gin.Default()

	route.POST("/subscription", subscriptionService.CreateSubscription)
	route.GET("/chat/:ChatID/subscription", subscriptionService.GetSubscriptionsByChatID)

	route.GET("/address", addressService.AddressGetList)
	route.GET("/address-count", addressService.AddressGetTotalCount)
	route.POST("/address-take", addressService.AddressTakeNext)
	route.POST("/address-take/:id", addressService.AddressTakeNext)
	route.GET("/address/:id", addressService.AddressGetByID)
	route.GET("/address-geocode/:lat/:lng/:accuracy", addressService.AddressGeocode)
	route.GET("/address-lookup/:id", addressService.AddressGeocodeByID)
	route.PUT("/address/:id", addressService.AddressUpdate)

	return route
}
