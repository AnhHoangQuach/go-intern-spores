package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/asaskevich/govalidator"
)

func (server *Server) signup(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ERROR(w, 400, err)
	}

	if govalidator.IsNull(user.Email) || govalidator.IsNull(user.Password) {
		utils.MessageError(w, 400, "Body is invalid")
	}

	user.Prepare()
	err = user.SignUp(server.DB, user.Email, user.Password, user.Phone, user.Address)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusCreated, "Please check code in your email")
}

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}