package app

import (
	"context"

	"github.com/ultraflame4/inuis/models"
	"gorm.io/gorm"
)


func CreateUser(username string, identity models.UserIdentity, grantRoles string, db *gorm.DB, ctx context.Context) (models.User, error) {
	user := models.User{}
	user.Username = username
	user.HID = GenerateHID()
	user.Roles = grantRoles
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[models.User](tx).Create(ctx, &user); err != nil {
			return err
		}
		// Attach identity
		identity.UserID = user.ID
		if err := gorm.G[models.UserIdentity](tx).Create(ctx, &identity); err != nil {
			return err
		}
		return nil
	});
	return user, err
}
