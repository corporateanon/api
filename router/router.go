package router

import (
	"github.com/gorilla/mux"
	"github.com/my1562/api/routes"
)

func NewRouter(
	subscriptionService *routes.SubscriptionService,
	addressService *routes.AddressService,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/subscription", subscriptionService.CreateSubscription).
		Methods("POST").
		Headers("Content-Type", "application/json")

	r.HandleFunc("/chat/{chatID}/subscription", subscriptionService.GetByChatID).
		Methods("GET")

	r.HandleFunc("/address/{id}", addressService.GetByID).
		Methods("GET")

	return r
}
