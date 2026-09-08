package models

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	// Hashed ID used for publicly retrieving (public) user info in the API.
	// Used to prevent ID iteration attacks. (Makes scraping information impossible by increasing IDs)
	HID string `json:"hid" gorm:"uniqueIndex"`
	// Metadata string delimited by semicolons
	// Never let this field be serialised!
	Roles string `json:"-"`
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


func (u *User) HasRole(role string) bool {
	if !strings.HasSuffix(role, ";") {
		role += ";"
	}
	return strings.Contains(u.Roles, role)
}

func (u *User) AddRole(role string) {
	if u.HasRole(role){
		return
	}
	if strings.HasSuffix(role, ";") {
		role += ";"
	}
	u.Roles += role
}

func (s *UserSession) ElevateNow(db *gorm.DB, ctx context.Context) {
	db.WithContext(ctx).Model(s).Update("elevated_at", time.Now())
}
