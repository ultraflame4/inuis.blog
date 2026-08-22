package main

import (
	"context"
	"log/slog"

	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/models"
	"github.com/teenyblog/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	slog.Info("Connecting to database...");
 	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  	if err != nil {
	    panic("failed to connect database")
  	}
    ctx := context.Background()
    _ = ctx;
    models.Migrate(db)
    slog.Info("Auto migrated database.")
    app.ResetSessions(db, ctx)
    slog.Info("Cleared all user sessions.")

	router := routers.InitRouter(db)
	router.Run(app.ServerAddress()) // listens on 0.0.0.0:8080 by default
}
