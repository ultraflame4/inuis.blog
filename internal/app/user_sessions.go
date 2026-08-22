package app

import (
	"context"
	"time"

	"github.com/teenyblog/models"
	"gorm.io/gorm"
)

// Deletes all sessions for a given user
func ClearUserSessions(db *gorm.DB, ctx context.Context, userID uint) error {
	_, err := gorm.G[models.UserSession](db).Where("user_id = ?", userID).Delete(ctx)
	return err
}

// Deletes all sessions
func ResetSessions(db *gorm.DB, ctx context.Context) error {
	_, err := gorm.G[models.UserSession](db).Where("1 = 1").Delete(ctx)
	return err
}


// Delete expired sessions
func ExpireSessions(db *gorm.DB, ctx context.Context) error {
	max_age := SessionMaxAge();
	now := time.Now()
	min_start := now.Add(-time.Duration(max_age) * time.Second)
	// Any session started before min_start is expired
	_, err := gorm.G[models.UserSession](db).Where("created_at < ?", min_start).Delete(ctx)
	return err
}

func CreateSession(db *gorm.DB, ctx context.Context, userID uint) (*models.UserSession, error) {
	token, err := RandString(128)
	if err != nil {
		return nil, err
	}
	session := models.UserSession{
		UserID: userID,
		Token:  token,
	}
	err = db.Create(&session).Error
	return &session, err
}

func DeleteSession(db *gorm.DB, ctx context.Context, sessionID uint) error {
	session := models.UserSession{
		ID: sessionID,
	}
	err := db.WithContext(ctx).Delete(&session).Error
	return err
}
