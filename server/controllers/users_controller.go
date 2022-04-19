package controllers

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"

	"heyui/server/models"
	"heyui/server/responses"
)

type UserController struct {
	DB *gorm.DB
}

func (u *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	user := models.User{}
	users, err := user.FetchUsers(u.DB, limit)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, toResponse(users))
}

func (u *UserController) GetUsersByFullName(w http.ResponseWriter, r *http.Request) {
	fullname := r.FormValue("fullname")

	user := models.User{}
	users, err := user.FindUserByFullName(u.DB, fullname)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, toResponse(users))
}

func (u *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["acct"]

	user := models.User{}
	userGotten, err := user.FindUserByID(u.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten.ToResponse())
}

func toResponse(users *[]models.User) []models.UserRep {
	usersRep := make([]models.UserRep, 0)
	for _, u := range *users {
		usersRep = append(usersRep, u.ToResponse())
	}
	return usersRep
}
