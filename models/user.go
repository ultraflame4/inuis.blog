package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
}

// User identity across different OIDCs
type UserIdentity struct {
	gorm.Model
	UserId uint
	User User
	Email *string `gorm:"index"`
	Provider string
	ProviderUserId string `gorm:"uniqueIndex"`
}
