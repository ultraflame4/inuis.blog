package main

import (
	"context"
	"fmt"

	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/models"
	"github.com/teenyblog/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Connecting to database...");
 	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  	if err != nil {
	    panic("failed to connect database")
  	}
    ctx := context.Background()
    _ = ctx;
    models.Migrate(db)
    fmt.Println("Auto migrated database.")
    app.ResetSessions(db, ctx)
    fmt.Println("Cleared all user sessions.")

	router := routers.InitRouter(db)
	router.Run(app.ServerAddress()) // listens on 0.0.0.0:8080 by default
}
