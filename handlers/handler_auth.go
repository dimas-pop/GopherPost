package handlers

import (
	"encoding/json"
	"fmt"
	"gopher-post/db"
	"gopher-post/utils"
	"net/http"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(s.DB, input.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	match := utils.CheckPasswordHash(input.Password, user.PasswordHash)
	if !match {
		http.Error(w, "Invalid check email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		fmt.Println("❌ ERROR JWT DETECTED:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"message": "Login Successful",
	})
}
