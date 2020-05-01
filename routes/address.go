package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/my1562/api/models"
	"github.com/my1562/api/notifier"
	"github.com/my1562/api/utils"
	"github.com/my1562/geocoder"
)

type AddressService struct {
	db       *gorm.DB
	geo      *geocoder.Geocoder
	notifier *notifier.Notifier
}

func NewAddressService(db *gorm.DB, geo *geocoder.Geocoder, notifier *notifier.Notifier) *AddressService {
	return &AddressService{
		db:       db,
		geo:      geo,
		notifier: notifier,
	}
}

type ShortGeocoderAddress struct {
	ID             uint32
	Address        string
	Building       string
	Street1562ID   uint32
	Street1562Name string
}

func FullToShortAddress(full *geocoder.FullAddress) (*ShortGeocoderAddress, error) {
	if full.Street1562 == nil {
		return nil, errors.New("Address without Street1562")
	}
	short := &ShortGeocoderAddress{
		ID:             full.Address.ID,
		Address:        formatAddress(full),
		Building:       full.Address.GetBuildingAsString(),
		Street1562ID:   full.Street1562.ID,
		Street1562Name: full.Street1562.Name,
	}
	return short, nil
}

type AddressUpdatePayload struct {
	CheckStatus    models.AddressArCheckStatus
	ServiceMessage string
	Hash           string
}

func (service *AddressService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}

	var payload AddressUpdatePayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}
	if payload.CheckStatus != models.AddressStatusNoWork && payload.CheckStatus != models.AddressStatusWork {
		utils.ErrorBadRequest(w, fmt.Sprintf("CheckStatus may be either '%s' or '%s'", models.AddressStatusNoWork, models.AddressStatusWork))
		return
	}

	address := &models.AddressAr{}
	err = service.db.
		Where(id).
		Preload("Subscriptions").
		First(address).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
	}
	if address.ID == 0 {
		utils.ErrorNotFound(w, "Not found")
		return
	}

	needsNotification := address.Hash != payload.Hash &&
		address.Subscriptions != nil &&
		len(address.Subscriptions) > 0

	address.CheckStatus = payload.CheckStatus
	address.ServiceMessage = payload.ServiceMessage
	address.Hash = payload.Hash
	address.CheckedAt = time.Now()
	err = service.db.Save(address).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
	}
	utils.Success(w, address)

	if needsNotification {
		chatIDs := make([]int64, len(address.Subscriptions))
		for i, subscription := range address.Subscriptions {
			chatIDs[i] = subscription.ChatID
		}
		geocoderAddress := service.geo.AddressByID(uint32(address.ID))
		service.notifier.NotifyServiceMessageChange(chatIDs, payload.ServiceMessage, formatAddress(geocoderAddress), address.CheckStatus)
	}
}

type GeocodeResponseAddress struct {
	ID            uint32
	Distance      float64
	AddressString string
}
type GeocodeResponse struct {
	Addresses []GeocodeResponseAddress
}

func formatAddress(addr *geocoder.FullAddress) string {
	bld := addr.Address.GetBuildingAsString()
	name := addr.StreetAR.NameRu
	t := addr.StreetAR.TypeRu
	return t + " " + name + " " + bld
}

func formatGeocodingResult(res *geocoder.ReverseGeocodingResult) string {
	return formatAddress(res.FullAddress)
}

func (service *AddressService) Geocode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error

	lat, err := strconv.ParseFloat(vars["lat"], 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}

	lng, err := strconv.ParseFloat(vars["lng"], 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}

	accuracy, err := strconv.ParseFloat(vars["accuracy"], 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}
	result := service.geo.ReverseGeocode(lat, lng, accuracy, 10)

	response := GeocodeResponse{Addresses: []GeocodeResponseAddress{}}
	for _, item := range result {
		response.Addresses = append(response.Addresses, GeocodeResponseAddress{
			ID:            item.FullAddress.Address.ID,
			Distance:      item.Distance,
			AddressString: formatGeocodingResult(item),
		})
	}
	utils.Success(w, response)
}

func (service *AddressService) GeocodeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}

	result := service.geo.AddressByID(uint32(id))
	if result == nil {
		utils.ErrorNotFound(w, "No such address")
		return
	}
	shortAddress, err := FullToShortAddress(result)
	if err != nil {
		utils.ErrorNotFound(w, err.Error())
		return
	}

	response := struct {
		Address *ShortGeocoderAddress
	}{Address: shortAddress}

	utils.Success(w, response)
}
