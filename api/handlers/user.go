package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"marketplace/internal/user"
	"net/http"
	"strconv"
)

var (
	UserStore = user.NewMemoryUserStore()
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var User user.User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user,_:=UserStore.CreateUser(&User)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var User *user.User
	User, err = UserStore.GetUser(id)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(User)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var updatedUser *user.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedUser.ID = id

	go func() {
		UserStore.UpdateUser(updatedUser)
	}()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("User with %v id was updated.", updatedUser.ID))
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = UserStore.DeleteUser(id)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("User with %v id was deleted.", id))
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Users, err := UserStore.GetAllUsers()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	json.NewEncoder(w).Encode(Users)
}