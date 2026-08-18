package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
}

type UserEmail struct {
	gorm.Model
	UserId uint
	User User
	Email string `gorm:"uniqueIndex"` // Primarily used to prevent creation of new accounts with different OIDC
}

// User identity across different OIDCs
type UserIdentity struct {

}
