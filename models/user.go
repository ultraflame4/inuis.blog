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
	Issuer string
	Subject string `gorm:"uniqueIndex"`
}

func FindOrCreateUserIdentity(){

}
