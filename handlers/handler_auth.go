package handlers

import (
	"encoding/json"
	"gopher-post/db"
	"gopher-post/utils"
	"log/slog"
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
		slog.WarnContext(r.Context(), "Invalid email or password", "error", err)
		JSONError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	match := utils.CheckPasswordHash(input.Password, user.PasswordHash)
	if !match {
		slog.WarnContext(r.Context(), "Invalid check email or password", "error", err)
		JSONError(w, "Invalid check email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error generating token", "error", err)
		JSONError(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Login succesfull")
	JSONSuccess(w, LoginResponse{Message: "login successful", Token: token}, http.StatusOK)
}
