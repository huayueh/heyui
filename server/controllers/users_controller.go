package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"heyui/server/auth"
	"heyui/utils/formaterror"
	"io/ioutil"
	"net/http"
	"strconv"

	"heyui/server/models"
	"heyui/server/responses"
)

type UserController struct {
	DB *gorm.DB
}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var token string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	inpwd := user.Pwd
	_, err = user.FindUserByID(u.DB, user.Acct)
	if err != nil {
		repLoginErr(user.Acct, w, http.StatusBadRequest, err)
		return
	}
	err = auth.VerifyPassword(user.Pwd, inpwd)
	if err != nil {
		repLoginErr(user.Acct, w, http.StatusUnauthorized, errors.New("incorrect account or password"))
		return
	}
	token, err = auth.CreateToken(user.Acct)

	responses.JSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (u *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("create")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	userCreated, err := user.SaveUser(u.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusCreated, userCreated.ToResponse())
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

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid := vars["acct"]

	_, err := user.DeleteAUser(u.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.SUCCESS(w, http.StatusOK)
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	vars := mux.Vars(r)
	user.Acct = vars["acct"]

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	updatedUser, err := user.UpdateUser(u.DB, user.Acct)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser.ToResponse())
}

func toResponse(users *[]models.User) []models.UserRep {
	usersRep := make([]models.UserRep, 0)
	for _, u := range *users {
		usersRep = append(usersRep, u.ToResponse())
	}
	return usersRep
}

func repLoginErr(acct string, w http.ResponseWriter, statusCode int, err error) {
	responses.ERROR(w, statusCode, err)
}
