package models

import (
	"fmt"
	"log"

	"github.com/AnhHoangQuach/go-intern-spores/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	confOption := config.GetConfigOption()
	if confOption == nil {
		fmt.Printf("Database is not have config")
	}
	database, err := gorm.Open(postgres.Open(confOption.PostgreDB), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("We are connected to the database")
	}

	database.AutoMigrate(User{}) //database migration

	DB = database
}