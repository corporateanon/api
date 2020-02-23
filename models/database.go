package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/my1562/api/config"
)

type Subscription struct {
	gorm.Model
	ChatID      int64 `gorm:"unique_index:idx_addr_chat"`
	AddressArID int64 `gorm:"unique_index:idx_addr_chat"`
}

type AddressArCheckStatus string

const (
	AddressStatusNoWork AddressArCheckStatus = "nowork"
	AddressStatusWork                        = "work"
	AddressStatusInit                        = "init"
)

type AddressAr struct {
	gorm.Model
	CheckStatus    AddressArCheckStatus
	ServiceMessage string
	TakenAt        time.Time
	CheckedAt      time.Time
	Subscriptions  []Subscription `gorm:"foreignkey:AddressArID"`
}

//NewDatabase creates database connection
func NewDatabase(conf *config.Config) (*gorm.DB, error) {
	db, err := gorm.
		Open(conf.DBDriver, conf.DBConnection)

	if err != nil {
		return nil, err
	}

	db = db.Debug()
	db = db.AutoMigrate(
		&Subscription{},
		&AddressAr{},
	)

	return db, nil
}
