package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}

// User identity across different OIDCs
type UserIdentity struct {
	gorm.Model
	UserID uint
	// Auto delete user identity when user is deleted
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email *string `gorm:"index"`
	Issuer string
	Subject string `gorm:"uniqueIndex"`
}

type UserSession struct {
	ID        uint `gorm:"primarykey"`
	Token string `gorm:"uniqueIndex"`
	UserID uint
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
}
