package models

import "github.com/jinzhu/gorm"

type Subscription struct {
	gorm.Model
	ChatID      int64 `gorm:"unique_index:idx_addr_chat"`
	AddressArID int64 `gorm:"unique_index:idx_addr_chat"`
}
