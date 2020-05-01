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

	r.HandleFunc("/address/geocode/{lat}/{lng}/{accuracy}", addressService.Geocode).
		Methods("GET")

	r.HandleFunc("/address/geocode/lookup/{id}", addressService.GeocodeByID).
		Methods("GET")

	r.HandleFunc("/address/{id}", addressService.GetByID).
		Methods("GET")

	r.HandleFunc("/address/{id}", addressService.Update).
		Methods("PUT").
		Headers("Content-Type", "application/json")

	return r
}
