package server

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	httpSwagger "github.com/swaggo/http-swagger"
	"heyui/server/auth"
	"heyui/server/controllers"
	"heyui/server/middlewares"

	_ "heyui/docs" // docs is generated by Swag CLI, you have to import it.
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

func (s *Server) initializePageRoutes() {
	s.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	s.Router.Path("/login").Handler(
		csrf.Protect(
			[]byte("538ce40f-2b40-4acd-bdf0-d1570b4168cf"),
			csrf.Secure(false),
		)(http.HandlerFunc(controllers.ShowLoginForm)),
	).Methods("GET")
}

func (s *Server) initializeApiRoutes() {
	uc := controllers.UserController{DB: s.DB}
	//Users routes
	s.Router.HandleFunc("/api/v1/users",
		middlewares.SetHeaders(
			middlewares.Auth(uc.GetUsersByFullName))).
		Queries("fullname", "{fullname}").
		Methods("GET")

	s.Router.HandleFunc("/api/v1/users",
		middlewares.SetHeaders(
			middlewares.Auth(uc.GetUsersByPage))).
		Queries("limit", "{limit}").
		Queries("page", "{page}").
		Methods("GET")

	s.Router.Path("/api/v1/auth/login/csrf").Handler(
		csrf.Protect(
			[]byte("538ce40f-2b40-4acd-bdf0-d1570b4168cf"),
			csrf.Secure(false),
		)(http.HandlerFunc(uc.Login)),
	).Methods("POST")

	s.Router.HandleFunc("/api/v1/auth/login",
		middlewares.SetHeaders(uc.Login)).
		Methods("POST")

	s.Router.HandleFunc("/api/v1/users",
		middlewares.SetHeaders(uc.CreateUser)).
		Methods("POST")

	s.Router.HandleFunc("/api/v1/users",
		middlewares.SetHeaders(
			middlewares.Auth(uc.GetUsers))).
		Methods("GET")

	s.Router.HandleFunc("/api/v1/users/{acct}",
		middlewares.SetHeaders(
			middlewares.Auth(uc.GetUser))).
		Methods("GET")

	s.Router.HandleFunc("/api/v1/users/{acct}",
		middlewares.SetHeaders(
			middlewares.Auth(uc.UpdateUser))).
		Methods("PUT")

	s.Router.HandleFunc("/api/v1/users/{acct}/fullname",
		middlewares.SetHeaders(
			middlewares.Auth(uc.UpdateFullname))).
		Methods("PUT")

	s.Router.HandleFunc("/api/v1/users/{acct}",
		middlewares.SetHeaders(
			middlewares.Auth(uc.DeleteUser))).
		Methods("DELETE")

	s.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

func (s *Server) Run(port string) {
	s.Router = mux.NewRouter()
	s.initializeApiRoutes()
	s.initializePageRoutes()
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

	go func() {
		addr := fmt.Sprintf(":%v", os.Getenv("WS_PORT"))
		route := mux.NewRouter()
		route.HandleFunc("/ws/v1/users/{acct}", controllers.WSconnection)
		http.ListenAndServe(addr, route)
	}()

	server := Server{DB: appdb}
	server.Run(os.Getenv("HTTP_PORT"))
}
