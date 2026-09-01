package app

import (
	"context"
	"github.com/ultraflame4/bytesblog/models"
	"gorm.io/gorm"
)


func CreateUser(username string, identity models.UserIdentity, db *gorm.DB, ctx context.Context) (models.User, error) {
	user := models.User{}
	user.HID = GenerateHID()
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
