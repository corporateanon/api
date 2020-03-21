package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/my1562/api/config"
)

//NewDatabase creates database connection
func NewDatabase(conf *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.
			Open(conf.DBDriver, conf.DBConnection)
		if err == nil {
			continue
		}
		fmt.Printf("Reconnecting to the database. Attempt %d", i)
		time.Sleep(5 * time.Second)
	}

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
