package app

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ultraflame4/bytesblog/models"
	"gorm.io/gorm"
)

func CreateUser(username string, identity models.UserIdentity, db *gorm.DB, ctx context.Context) (models.User, error) {
	user := models.User{}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[models.User](tx).Create(ctx, &user); err != nil {
			return err
		}
		// Attach identity
		if err := gorm.G[models.UserIdentity](tx).Create(ctx, &identity); err != nil {
			return err
		}

		hid_s := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), user.ID)
		hash := sha512.Sum512([]byte(hid_s))
		hid := fmt.Sprintf("%d%s", user.ID, base64.RawStdEncoding.EncodeToString(hash[:]))
		user.HID = &hid

		return nil
	});
	return user, err
}
