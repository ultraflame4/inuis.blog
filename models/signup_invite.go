package models

import (
	"time"
)

type SignupInviteCode struct {
	ID        uint `gorm:"primarykey"`
	Code      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	// Longest Validity period for invite code is 7 days
	Validity  *time.Duration
	SingleUse  bool `gorm:"default:true"`
	// Additional metadata for specific use cases
	Metadata  *string
}

func (s *SignupInviteCode) IsExpired() bool {
	// Check if the invite code has expired after 7 days
	if time.Since(s.CreatedAt) > time.Hour * 24 * 7 {
		return true
	}
	// If validity is not set, check if code is past its validity period
	if s.Validity != nil {
		return time.Since(s.CreatedAt) > *s.Validity
	}
	return false
}

