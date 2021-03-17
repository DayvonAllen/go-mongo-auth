package app

import (
	"example.com/app/repo"
	"example.com/app/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {

	router := mux.NewRouter()

	// wiring everything up
	ch := Handlers{userService: services.NewUserService(repo.NewUserRepoImpl()),
		authService: services.NewAuthService(repo.NewAuthRepoImpl())}

	router.HandleFunc("/users", ch.getAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users", ch.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/login", ch.Login).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", ch.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", ch.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", ch.DeleteByID).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}


