package models

import (
	"context"
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
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

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
	// If not null, time of which session was elevated.
	// Elevated session only lasts for 30seconds.
	// Elevated sessions allow user to perform more risky actions.
	// To elevate a session, use is required to sign in again.
	ElevatedAt *time.Time `gorm:"default:null"`
}

func (s *UserSession) IsElevated() bool {
	if s.ElevatedAt == nil {
		return false
	}
	return time.Since(*s.ElevatedAt) < 30*time.Second
}

func (s *UserSession) ElevateNow(db *gorm.DB, ctx context.Context) {
	db.WithContext(ctx).Model(s).Update("elevated_at", time.Now())
}
