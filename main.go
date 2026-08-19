package main

import (
	"context"
	"fmt"

	"github.com/teenyblog/models"
	"github.com/teenyblog/routers"
	"github.com/teenyblog/utils"
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
    db.AutoMigrate(&models.User{})
    fmt.Println("Auto migrated database.")

	router := routers.InitRouter()
	router.Run(utils.GetServerAddress()) // listens on 0.0.0.0:8080 by default
}
