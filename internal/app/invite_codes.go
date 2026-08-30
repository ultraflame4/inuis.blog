package app

import (
	"context"

	"github.com/ultraflame4/bytesblog/models"
	"gorm.io/gorm"
)

// Returns true if the invite code is valid, false otherwise.
//
// This will also delete the invite code from the database (if expired or single use).
func ConsumeInviteCode(db *gorm.DB, ctx context.Context, code string) bool {
	var inviteCode models.SignupInviteCode
	if err := db.WithContext(ctx).Where("code = ?", code).First(&inviteCode).Error; err != nil {
		expired := inviteCode.IsExpired();
		// Delete the invite code from the database if expired or is single use.
		if expired || inviteCode.SingleUse {
			db.WithContext(ctx).Delete(inviteCode)
		}
		return !expired
	}
	return false
}