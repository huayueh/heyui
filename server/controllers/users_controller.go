package controllers

import (
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

func toResponse(users *[]models.User) []models.UserRep {
	usersRep := make([]models.UserRep, 0)
	for _, u := range *users {
		usersRep = append(usersRep, u.ToResponse())
	}
	return usersRep
}
