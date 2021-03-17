package app

import (
	"encoding/json"
	"example.com/app/domain"
	"example.com/app/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type Handlers struct {
	userService services.UserService
	authService services.AuthService
}

func (ch *Handlers) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ch.userService.GetAllUsers()

	if err != nil {
		log.Panicf("error: %v", err)
	}

	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (ch *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var users domain.User
	err := json.NewDecoder(r.Body).Decode(&users)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = ch.userService.CreateUser(&users)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (ch *Handlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(name["id"])

	user, err := ch.userService.GetUserByID(id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (ch *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	name := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(name["id"])

	u, _ := ch.userService.UpdateUser(id, &user)
	err = json.NewEncoder(w).Encode(u)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (ch *Handlers) DeleteByID(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(name["id"])

	err := ch.userService.DeleteByID(id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (ch *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var credential domain.LoginDetails
	var auth domain.Authentication
	err := json.NewDecoder(r.Body).Decode(&credential)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	u, token, err := ch.authService.Login(credential.Email, credential.Password)

	if err != nil {
		http.Error(w, "invalid credential", http.StatusUnauthorized)
		return
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer " + token + "|")...)
	t, err := auth.SignToken(signedToken)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	signedToken = append(signedToken, t...)

	c := http.Cookie{
		Name:  "session",
		Value: string(signedToken),
	}

	http.SetCookie(w, &c)

	err = json.NewEncoder(w).Encode(&u)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}