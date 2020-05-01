package routes

import (
	"github.com/jinzhu/gorm"
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
