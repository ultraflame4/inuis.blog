package app

import (
	"context"

	"github.com/ultraflame4/inuis/models"
	"gorm.io/gorm"
)

// WhereUserWithRole finds a user with the given roles.
//
// Returns the partial query interface for chaining.
func WhereUserWithRole(role string, db *gorm.DB, ctx context.Context) gorm.ChainInterface[models.User] {
	return gorm.G[models.User](db).Where("roles LIKE ?", "%"+role+";%")
}

// FindUserWithRole finds a user with the given roles.
//
// Returns nil if no user is found.
func FindUserWithRole(role string, db *gorm.DB, ctx context.Context) *models.User {
	user, err := WhereUserWithRole(role, db, ctx).First(ctx)
	if err != nil {
		return nil
	}
	return &user
}

func CountUsersWithRole(role string, db *gorm.DB, ctx context.Context) int64 {
	count, err := WhereUserWithRole(role, db, ctx).Count(ctx, "ID")
	if err != nil {
		panic("Unexpected error while CountUsersWithRole(" + role + "):" + err.Error())
	}
	return count
}

func FindUserByHid(hid string, db *gorm.DB, ctx context.Context) (models.User, bool) {
	user, err := gorm.G[models.User](db).Where("h_id = ?", hid).First(ctx)
	return user, err == nil
}
func FindUserIdentity(subject string, issuer string, db *gorm.DB, ctx context.Context) (models.UserIdentity, error) {
	identity, err := gorm.G[models.UserIdentity](db).
		Preload("User", nil).
		Where("subject = ? AND issuer = ?", subject, issuer).
		First(ctx)
	return identity, err
}
