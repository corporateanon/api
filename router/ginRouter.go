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
	route.GET("/address/:id", addressService.AddressGetByID)

	// r.HandleFunc("/address/geocode/{lat}/{lng}/{accuracy}", addressService.Geocode).
	// 	Methods("GET")

	// r.HandleFunc("/address/geocode/lookup/{id}", addressService.GeocodeByID).
	// 	Methods("GET")

	// r.HandleFunc("/address/{id}", addressService.GetByID).
	// 	Methods("GET")

	// r.HandleFunc("/address/{id}", addressService.Update).
	// 	Methods("PUT").
	// 	Headers("Content-Type", "application/json")

	return route
}