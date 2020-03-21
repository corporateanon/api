package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AddressArCheckStatus string

const (
	AddressStatusNoWork AddressArCheckStatus = "nowork"
	AddressStatusWork                        = "work"
	AddressStatusInit                        = "init"
)

type AddressAr struct {
	gorm.Model
	CheckStatus    AddressArCheckStatus
	ServiceMessage string `gorm:"size:2048"`
	Hash           string
	TakenAt        time.Time
	CheckedAt      time.Time
	Subscriptions  []Subscription `gorm:"foreignkey:AddressArID"`
}
