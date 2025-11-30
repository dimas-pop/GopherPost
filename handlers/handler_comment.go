package handlers

import (
	"encoding/json"
	"gopher-post/db"
	"gopher-post/middleware"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// GetCommentHandler godoc
// @Summary      Dapatkan komentar
// @Description  Mengambil semua komentar berdasarkan ID Postingan
// @Tags         comments
// @Produce      json
// @Param        id path   string  true  "ID Postingan (UUID)"
// @Success      200    {array}  models.Comment
// @Failure	     400	{object} handlers.ErrorResponse
// @Router       /posts/{id}/comments [get]
func (s *Server) GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	result, err := db.GetCommentByPostID(s.DB, postID)
	if err != nil {
		JSONError(w, "Failed get comment", http.StatusBadRequest)
		return
	}

	JSONSuccess(w, &result, http.StatusOK)
}

// CreateCommentHandler godoc
// @Summary      Kirim komentar
// @Description  Memberikan komentar pada postingan tertentu
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id path   string  true  "ID Postingan (UUID)"
// @Param        request body   handlers.CreateCommentInput true "Isi Komentar"
// @Success      201  {object}  handlers.SuccessResponse
// @Failure	     400  {object}  handlers.ErrorResponse
// @Failure	     500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/posts/{id}/comments [post]
func (s *Server) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateCommentInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		JSONError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID := vars["id"]
	userID := r.Context().Value(middleware.UserIDKey).(string)

	err = db.CreateCommentInDB(s.DB, input.Content, userID, postID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed create comment", "error", err)
		JSONError(w, "Failed create comment", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Comment created")
	JSONSuccess(w, SuccessResponse{Message: "comment created"}, http.StatusCreated)
}

// DeleteCommentHandler godoc
// @Summary      Hapus komentar
// @Description  Menghapus comment permanen berdasarkan ID
// @Tags         comments
// @Produce      json
// @Param        id   path      string  true  "ID Komentar (UUID)"
// @Security     BearerAuth
// @Success      200  {object}  handlers.SuccessResponse
// @Failure	     404  {object}  handlers.ErrorResponse
// @Failure	     403  {object}  handlers.ErrorResponse
// @Failure	     500  {object}  handlers.ErrorResponse
// @Router       /api/comments/{id} [delete]
func (s *Server) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]

	currentUserID := r.Context().Value(middleware.UserIDKey).(string)

	userID, err := db.GetCommentOwnerID(s.DB, commentID)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed found user_id", "error", err)
		JSONError(w, "Failed get user_id", http.StatusNotFound)
		return
	}

	if currentUserID != userID {
		slog.WarnContext(r.Context(), "Invalid user_id")
		JSONError(w, "Invalid user_id", http.StatusForbidden)
		return
	}

	err = db.DeleteCommentByID(s.DB, commentID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed delete comment", "error", err)
		JSONError(w, "Failed delete comment", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Comment deleted")
	JSONSuccess(w, SuccessResponse{Message: "comment deleted"}, http.StatusOK)
}
