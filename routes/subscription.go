package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
	"github.com/my1562/geocoder"
)

type SubscriptionService struct {
	db  *gorm.DB
	geo *geocoder.Geocoder
}

func NewSubscriptionService(db *gorm.DB, geo *geocoder.Geocoder) *SubscriptionService {
	return &SubscriptionService{
		db:  db,
		geo: geo,
	}
}

type SubscriptionRequestPost struct {
	ChatID      int64
	AddressArID int64
}

func (service *SubscriptionService) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var err error
	var payload SubscriptionRequestPost
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}
	if payload.ChatID <= 0 {
		utils.ErrorBadRequest(w, "ChatID <= 0")
		return
	}
	if payload.AddressArID <= 0 {
		utils.ErrorBadRequest(w, "AddressArID <= 0")
		return
	}

	addressFromGeocoder := service.geo.AddressByID(uint32(payload.AddressArID))
	if addressFromGeocoder == nil {
		utils.ErrorBadRequest(w, "AddressArID not found")
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
		utils.Success(w, existingSubscription)
		tx.Commit()
		return
	}

	tx.Save(subscription)

	err = tx.Commit().Error
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}
	utils.Success(w, subscription)
}

func (service *SubscriptionService) GetByChatID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID, err := strconv.ParseUint(vars["chatID"], 10, 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}
	if chatID <= 0 {
		utils.ErrorBadRequest(w, "chatID <= 0")
		return
	}

	var subscriptions []models.Subscription
	service.db.
		Where(models.Subscription{ChatID: int64(chatID)}).
		Find(&subscriptions)

	utils.Success(w, subscriptions)
}
