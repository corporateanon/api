package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/my1562/api/models"
	"github.com/my1562/api/utils"
)

type AddressService struct {
	db *gorm.DB
}

func NewAddressService(db *gorm.DB) *AddressService {
	return &AddressService{
		db: db,
	}
}

func (service *AddressService) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorBadRequest(w, err.Error())
		return
	}
	address := &models.AddressAr{}

	service.db.
		Preload("Subscriptions").
		Where(id).
		First(address)
	if address.ID == 0 {
		utils.ErrorNotFound(w, "Address not found")
		return
	}

	utils.Success(w, address)
}

func (service *AddressService) GetTotalCount(w http.ResponseWriter, r *http.Request) {
	address := &models.AddressAr{}
	var cnt int64 = 0
	err := service.db.Model(address).Count(&cnt).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
	}
	utils.Success(w, cnt)
}

// TakeNext takes the oldest taken address and updates its TakenAt field to the present moment
func (service *AddressService) TakeNext(w http.ResponseWriter, r *http.Request) {
	address := &models.AddressAr{}
	err := service.db.Order("taken_at ASC").First(address).Error
	if err != nil {
		if err.Error() == "record not found" {
			utils.ErrorNotFound(w, "No addresses in database")
			return
		} else {
			utils.ErrorInternal(w, err.Error())
			return
		}
	}
	if address.ID == 0 {
		utils.ErrorNotFound(w, "No addresses in database")
		return
	}

	address.TakenAt = time.Now()
	err = service.db.Save(address).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
	}
	utils.Success(w, address)
}

type AddressUpdatePayload struct {
	CheckStatus    models.AddressArCheckStatus
	ServiceMessage string
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
	err = service.db.Where(id).First(address).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
	}
	if address.ID == 0 {
		utils.ErrorNotFound(w, "Not found")
		return
	}
	address.CheckStatus = payload.CheckStatus
	address.ServiceMessage = payload.ServiceMessage
	address.CheckedAt = time.Now()
	err = service.db.Save(address).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
	}
	utils.Success(w, address)
}
