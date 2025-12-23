package handlers

import (
	"encoding/json"
	"gopher-post/db"
	"gopher-post/middleware"
	"gopher-post/utils"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// GetUserAllHandler godoc
// @Summary      List all users
// @Description  Retrieves a list of all registered users
// @Tags         users
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /api/users [get]
func (s *Server) GetUserAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUserAll(s.DB)
	if err != nil {
		utils.JSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, &users, http.StatusOK)
}

// GetUserByIDHandler godoc
// @Summary      Lihat profil user
// @Description  Mengambil data detail user berdasarkan ID (UUID)
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "ID User (UUID)"
// @Success      200  {object}  models.User
// @Failure      404  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /api/users/{id} [get]
func (s *Server) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := db.GetUserByID(s.DB, id)
	if err != nil {
		utils.JSONError(w, "User not found or database error", http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, &user, http.StatusOK)
}

// CreateUserHandler godoc
// @Summary      Daftar user baru
// @Description  Mendaftarkan akun baru ke sistem
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body handlers.RegisterInput true "Data User"
// @Success      201  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      409  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /register [post]
func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, "Bad Request Input", http.StatusBadRequest)
		return
	}

	exists, err := db.CheckEmailExists(s.DB, input.Email)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed check email in DB",
			"error", err,
			"email", input.Email,
		)
		utils.JSONError(w, "Failed database check", http.StatusInternalServerError)
		return
	}

	if exists {
		utils.JSONError(w, "Email already in use", http.StatusConflict)
		return
	}

	password_hash, err := utils.HashPassword(input.Password)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed hashing password", "error", err)
		utils.JSONError(w, "Failed hash password", http.StatusInternalServerError)
		return
	}

	err = db.CreateUserInDB(s.DB, input.Name, input.Email, password_hash)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed create user in DB",
			"error", err,
			"name", input.Name,
			"email", input.Email,
		)
		utils.JSONError(w, "Failed create user", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "User created successfully",
		"name", input.Name,
		"email", input.Email,
	)
	utils.JSONSuccess(w, utils.SuccessResponse{Message: "user created"}, http.StatusCreated)
}

// UpdateUserHandler godoc
// @Summary      Update user profile
// @Description  Updates the name and email of a user identified by ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Param        request body handlers.UpdateUserInput true "Updated user data"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /api/users/{id} [put]
func (s *Server) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input UpdateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	currentUserID := r.Context().Value(middleware.UserIDKey).(string)
	if currentUserID != id {
		slog.WarnContext(r.Context(), "Invalid user",
			"attempt_by_user_id", currentUserID,
			"target_owner_id", id,
		)
		utils.JSONError(w, "Invalid user", http.StatusForbidden)
		return
	}

	err = db.UpdateUserByID(s.DB, input.Name, input.Email, id)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed update user in DB",
			"error", err,
			"user_id", id)
		utils.JSONError(w, "Failed update user", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "User updated successfully",
		"user_id", id,
	)
	utils.JSONSuccess(w, utils.SuccessResponse{Message: "user updated"}, http.StatusOK)
}

// DeleteUserHandler godoc
// @Summary      Delete user
// @Description  Deletes a user identified by ID.
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /api/users/{id} [delete]
func (s *Server) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	currentUserID := r.Context().Value(middleware.UserIDKey).(string)
	if currentUserID != id {
		slog.WarnContext(r.Context(), "Invalid user",
			"attempt_by_user_id", currentUserID,
			"target_owner_id", id,
		)
		utils.JSONError(w, "Invalid user", http.StatusForbidden)
		return
	}

	err := db.DeleteUserByID(s.DB, id)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed delete user in DB",
			"error", err,
			"user_id", id,
		)
		utils.JSONError(w, "Failed delete user", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "User deleted successfully",
		"user_id", id,
	)
	utils.JSONSuccess(w, utils.SuccessResponse{Message: "user deleted"}, http.StatusOK)
}
