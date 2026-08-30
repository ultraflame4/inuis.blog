package app

import (
	"context"
	"crypto/rand"
	"log/slog"
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
		expired := inviteCode.IsExpired()
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
func GenerateInviteCode(db *gorm.DB, ctx context.Context, validity *time.Duration, singleUse bool, roles []string) (string, error) {
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
	inviteCode.SetRoles(roles)
	if err := db.WithContext(ctx).Create(&inviteCode).Error; err != nil {
		return "", err
	}
	return inviteCode.Code, nil
}

// Purge expired invite codes
func PurgeExpiredInviteCodes(db *gorm.DB, ctx context.Context) error {
	inviteCodes, err := gorm.G[models.SignupInviteCode](db).Find(ctx)
	if err != nil {
		return err
	}
	count := 0
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, inviteCode := range inviteCodes {
			if inviteCode.IsExpired() {
				if err := tx.WithContext(ctx).Delete(&inviteCode).Error; err != nil {
					return err
				}
				count++
			}
		}
		return nil
	})
	if err != nil {
		slog.Error("Error while purging expired invite codes", "error", err)
		return err
	}
	slog.Info("Purged expired invite codes", "count", count)
	return nil
}


// Generate admin invite codes
func GenerateAdminInviteCode(db *gorm.DB, ctx context.Context) (string, error) {
	validity := time.Hour * 1
	
	return GenerateInviteCode(db, ctx, &validity, true, []string{"admin"})
}

func GetAdminInviteCode(db *gorm.DB, ctx context.Context) (string, error) {
	PurgeExpiredInviteCodes(db, ctx)
	role := "admin"
	inviteCode, err := gorm.G[models.SignupInviteCode](db).Where("grant_roles LIKE ?", "%"+role+";%").First(ctx)
	if inviteCode.IsExpired() || err != nil {
		slog.Info("No admin invite codes found. Generating new one... (Note that admin invite codes are only valid for 1 hour)")
		code, err := GenerateAdminInviteCode(db, ctx)
		if err != nil {
			return "", err
		}
		return code, nil
	}
	return inviteCode.Code, nil
}
