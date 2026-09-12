package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/ultraflame4/inuis/internal/app"
	"github.com/ultraflame4/inuis/models"
	"github.com/ultraflame4/inuis/routers"
)

//go:embed all:embed
var frontendFS embed.FS


func main() {
	app_config := app.GetConfig()
	
	slog.Info("Connecting to database...")
	
	db, err := app_config.ConnectDatabase()
	if err != nil {
		slog.Error("Failed to open connection to database.", "error",err)
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


	if app.CountUsersWithRole("admin", db, ctx) == 0 {
		slog.Info("No admin users found! Generating invite code for admin user!")
		inviteCode, err := app.GetAdminInviteCode(db, ctx)
		if err != nil {
			slog.Error("Failed to generate [admin] invite code", "error", err)
		} else {
			slog.Info("Generated [admin] invite code", "code", inviteCode)
			adminSetupURL := app_config.PublicServerAddr() + "/signup?invite=" + inviteCode
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
	router.Run(app_config.ServerAddress()) // listens on 0.0.0.0:8080 by default
}
