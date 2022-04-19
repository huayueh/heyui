package controllers

import (
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
)

type PageController struct {
	DB *gorm.DB
}

func ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("wsclient.html"))
	t.Execute(w, map[string]interface{}{})
}
