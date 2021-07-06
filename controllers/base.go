package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AnhHoangQuach/go-intern-spores/config"
	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize() error{
	var err error
	confOption := config.GetConfigOption()
	if confOption == nil {
		return nil
	}
	server.DB, err = gorm.Open(postgres.Open(confOption.PostgreDB), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	} else {
		fmt.Printf("We are connected to the database")
	}

	server.DB.Debug().AutoMigrate(&models.User{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
	return nil
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}