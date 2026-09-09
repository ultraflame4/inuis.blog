package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/ultraflame4/inuis/internal/app"
	"github.com/ultraflame4/inuis/models"
	"github.com/ultraflame4/inuis/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed embed/*
var frontendFS embed.FS


func main() {
	slog.Info("Connecting to database...")
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	ctx := context.Background()
	_ = ctx
	models.Migrate(db)
	slog.Info("Auto migrated database.")
	app.ExpireSessions(db, ctx)
	slog.Info("Cleared all user sessions.")
	app.PurgeExpiredInviteCodes(db, ctx)
	slog.Info("Purged expired invite codes.")

	aconfig := app.GetConfig()
	if app.CountUsersWithRole("admin", db, ctx) == 0 {
		slog.Info("No admin users found! Generating invite code for admin user!")
		inviteCode, err := app.GetAdminInviteCode(db, ctx)
		if err != nil {
			slog.Error("Failed to generate [admin] invite code", "error", err)
		} else {
			slog.Info("Generated [admin] invite code", "code", inviteCode)
			adminSetupURL := aconfig.PublicServerAddr() + "/signup?invite=" + inviteCode
			fmt.Printf(`
╔═══════════════════════════════════════════════════════════════════════════════════════╗
║                                      Admin Setup                                      ║
╚═══════════════════════════════════════════════════════════════════════════════════════╝

	Create the initial admin user by opening this URL:

	http://%s

	Follow the instructions to complete the setup.

╔═══════════════════════════════════════════════════════════════════════════════════════╗
╚═══════════════════════════════════════════════════════════════════════════════════════╝
            `, adminSetupURL)
		}
	}

	router := routers.InitRouter(db, frontendFS)
	router.Run(aconfig.ServerAddress()) // listens on 0.0.0.0:8080 by default
}
