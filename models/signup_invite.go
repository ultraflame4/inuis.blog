package models

import (
	"strings"
	"time"
)

type SignupInviteCode struct {
	ID        uint `gorm:"primarykey"`
	Code      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	// Longest Validity period for invite code is 7 days
	Validity  *time.Duration
	SingleUse  bool `gorm:"default:true"`
	// Roles to grant to the user when they sign up. Note that this should be a semicolon-separated list of roles.
	// Needs to also be suffixed with ";"
	GrantRoles  *string
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

// Roles to grant to the user when they sign up. Don't need to be a semicolon-separated list.
func (s *SignupInviteCode) SetRoles(roles []string) {
	if len(roles) == 0 {
		s.GrantRoles = nil
		return
	}
	rs := strings.Join(roles, ";") + ";"
	s.GrantRoles = &rs
}