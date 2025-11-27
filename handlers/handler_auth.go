package handlers

import (
	"encoding/json"
	"gopher-post/db"
	"gopher-post/utils"
	"net/http"
)

// LoginHandler godoc
// @Summary      Masuk ke aplikasi
// @Description  Tukar email dan password dengan Token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body handlers.LoginInput true "Kredensial Login"
// @Success      200  {object}  handlers.LoginResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /login [post]
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		JSONError(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(s.DB, input.Email)
	if err != nil {
		JSONError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	match := utils.CheckPasswordHash(input.Password, user.PasswordHash)
	if !match {
		JSONError(w, "Invalid check email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		JSONError(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, LoginResponse{Message: "login successful", Token: token}, http.StatusOK)
}
