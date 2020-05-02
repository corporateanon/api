package routes

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/my1562/api/models"
	"github.com/my1562/api/notifier"
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

func formatAddress(addr *geocoder.FullAddress) string {
	bld := addr.Address.GetBuildingAsString()
	name := addr.StreetAR.NameRu
	t := addr.StreetAR.TypeRu
	return t + " " + name + " " + bld
}

func formatGeocodingResult(res *geocoder.ReverseGeocodingResult) string {
	return formatAddress(res.FullAddress)
}
