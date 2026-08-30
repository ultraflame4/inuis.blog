package app

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

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

// Generates and inserts a new invite code into the database.
// 
// Returns the generated invite code, or an empty string if the operation fails.
func GenerateInviteCode(db *gorm.DB, ctx context.Context, validity *time.Duration, singleUse bool) (string, error){
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 16)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}

	var inviteCode models.SignupInviteCode
	inviteCode.Code = string(code)
	inviteCode.Validity = validity
	inviteCode.SingleUse = singleUse
	if err := db.WithContext(ctx).Create(&inviteCode).Error; err != nil {
		return "", err
	}
	return inviteCode.Code, nil
}