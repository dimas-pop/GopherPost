package handlers

import (
	"encoding/json"
	"fmt"
	"gopher-post/db"
	"gopher-post/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetUserAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUserAll(s.DB)
	if err != nil {
		http.Error(w, "Failed get user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&users)
}

func (s *Server) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user, err := db.GetUserByID(s.DB, id)
	if err != nil {
		http.Error(w, "Failed get user by ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	exists, err := db.CheckEmailExists(s.DB, input.Email)
	if err != nil {
		http.Error(w, "Failed database check", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	password_hash, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Failed hash password", http.StatusInternalServerError)
		return
	}
	err = db.CreateUserInDB(s.DB, input.Name, input.Email, password_hash)
	if err != nil {
		http.Error(w, "Failed create user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func (s *Server) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	err = db.UpdateUserByID(s.DB, input.Name, input.Email, id)
	if err != nil {
		http.Error(w, "Failed update user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
}

func (s *Server) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.DeleteUserByID(s.DB, id)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Failed delete user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User Deleted"})
}
