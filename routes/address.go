package routes

import (
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
		Where(id).
		First(address)

	utils.Success(w, address)
}

// TakeNext takes the oldest taken address and updates its TakenAt field to the present moment
func (service *AddressService) TakeNext(w http.ResponseWriter, r *http.Request) {
	address := &models.AddressAr{}
	err := service.db.Order("taken_at ASC").First(address).Error
	if err != nil {
		utils.ErrorInternal(w, err.Error())
		return
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
