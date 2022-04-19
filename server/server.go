package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"heyui/server/auth"
	"heyui/server/controllers"
	"heyui/server/middlewares"

	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/joho/godotenv"
	"heyui/server/db"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) initializeApiRoutes() {
	uc := controllers.UserController{DB: s.DB}
	//Users routes

	s.Router.HandleFunc("/api/v1/users",
		middlewares.SetHeaders(
			middlewares.Auth(uc.GetUsers))).
		Methods("GET")
}

func (s *Server) Run(port string) {
	s.Router = mux.NewRouter()
	s.initializeApiRoutes()
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Listening to port %v\n", port)
	http.ListenAndServe(addr, s.Router)
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Print(".env file found")
	}
}

func Start() {
	var err error
	err = godotenv.Load()
	if err != nil {
		panic(".env required to launch the RESTful service")
	}

	appdb := db.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	auth.Initialize()

	server := Server{DB: appdb}
	server.Run(os.Getenv("HTTP_PORT"))
}
