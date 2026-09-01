package app

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/ultraflame4/bytesblog/models"
	"gorm.io/gorm"
)

func CreateUser(username string, identity models.UserIdentity, db *gorm.DB, ctx context.Context) (models.User, error) {
	user := models.User{}
	uuid_ := uuid.New()
	user.HID = base64.RawStdEncoding.EncodeToString(uuid_[:])
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[models.User](tx).Create(ctx, &user); err != nil {
			return err
		}
		// Attach identity
		if err := gorm.G[models.UserIdentity](tx).Create(ctx, &identity); err != nil {
			return err
		}


		return nil
	});
	return user, err
}
