package handlers

import (
	"encoding/json"
	"gopher-post/db"
	"gopher-post/middleware"
	"gopher-post/utils"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllPostsHandler godoc
// @Summary      Melihat semua postingan
// @Description  Mengambil daftar postingan dengan pagination
// @Tags         posts
// @Produce      json
// @Param        page  query    int     false  "Nomer Halaman"
// @Param        limit query    int     false  "Jumlah Data per Halaman"
// @Success      200   {array}  models.Post
// @Failure      500   {object} handlers.ErrorResponse
// @Router       /posts [get]
func (s *Server) GetPostAllHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	if limit < 1 {
		limit = 10
	}

	offSet := (limit * page) - limit

	posts, err := db.GetPostAll(s.DB, limit, offSet)
	if err != nil {
		utils.JSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, &posts, http.StatusOK)
}

// GetPostByIDHandler godoc
// @Summary      Melihat satu postingan
// @Description  Mengambil detail post berdasarkan ID
// @Tags         posts
// @Produce      json
// @Param        id   path      string  true  "ID Postingan (UUID)"
// @Success      200   {object}  models.Post
// @Failure      400   {object}  handlers.ErrorResponse
// @Failure      404   {object}  handlers.ErrorResponse
// @Failure      500   {object}  handlers.ErrorResponse
// @Router       /posts/{id} [get]
func (s *Server) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	post, err := db.GetPostByID(s.DB, id)
	if err != nil {
		utils.JSONError(w, "Post not found or database error", http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, &post, http.StatusOK)
}

// CreatePostHandler godoc
// @Summary      Membuat postingan baru
// @Description  Membuat post dengan judul dan konten. Butuh token JWT.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request body handlers.CreatePostInput true "Data Post"
// @Success      201  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/posts [post]
func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost CreatePostInput
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		utils.JSONError(w, "Bad request input", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		slog.WarnContext(r.Context(), "Auth Context missing UserID")
		utils.JSONError(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	err = db.CreatePostInDB(s.DB, newPost.Title, newPost.Content, userID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed create post in DB",
			"error", err,
			"user_id", userID,
		)
		utils.JSONError(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Post created successfully",
		"user_id", userID,
	)
	utils.JSONSuccess(w, utils.SuccessResponse{Message: "post created"}, http.StatusCreated)
}

// UpdatePostHandler godoc
// @Summary      Edit postingan
// @Description  Mengubah judul atau konten post berdasarkan ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID Postingan (UUID)"
// @Param        request body   handlers.UpdatePostInput true "Data Update"
// @Security     BearerAuth
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /api/posts/{id} [put]
func (s *Server) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	currentUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || currentUserID == "" {
		slog.ErrorContext(r.Context(), "Auth Context missing UserID")
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
	}

	ownerID, err := db.GetPostOwnerID(s.DB, postID)
	if err != nil {
		slog.WarnContext(r.Context(), "Update failed: Not found post",
			"error", err,
			"post_id", postID,
		)
		utils.JSONError(w, "Failed get post", http.StatusNotFound)
		return
	}

	if currentUserID != ownerID {
		slog.WarnContext(r.Context(), "Update failed: Forbidden access",
			"post_id", postID,
			"attempt_by_user_id", currentUserID,
			"target_owner_id", ownerID,
		)
		utils.JSONError(w, "You are not allowed to update this post", http.StatusForbidden)
		return
	}

	var input UpdatePostInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, "Bad request input", http.StatusBadRequest)
		return
	}

	err = db.UpdatePostByID(s.DB, input.Title, input.Content, postID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to update post in DB",
			"error", err,
			"post_id", postID,
		)
		utils.JSONError(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Post updated successfully",
		"post_id", postID,
		"user_id", currentUserID,
	)
	utils.JSONSuccess(w, utils.SuccessResponse{Message: "post updated"}, http.StatusOK)
}

// DeletePostHandler godoc
// @Summary      Hapus postingan
// @Description  Menghapus post permanen berdasarkan ID
// @Tags         posts
// @Produce      json
// @Param        id   path      string  true  "ID Postingan (UUID)"
// @Security     BearerAuth
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /api/posts/{id} [delete]
func (s *Server) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	currentUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || currentUserID == "" {
		slog.ErrorContext(r.Context(), "Auth Context missing UserID")
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
	}

	ownerID, err := db.GetPostOwnerID(s.DB, postID)
	if err != nil {
		slog.WarnContext(r.Context(), "Not found post",
			"error", err,
			"post_id", postID,
		)
		utils.JSONError(w, "Failed to get post", http.StatusNotFound)
		return
	}

	if currentUserID != ownerID {
		slog.WarnContext(r.Context(), "Invalid user_id",
			"post_id", postID,
			"attempt_by_user_id", currentUserID,
			"target_owner_id", ownerID,
		)
		utils.JSONError(w, "Invalid user_id", http.StatusForbidden)
		return
	}

	err = db.DeletePostByID(s.DB, postID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed delete post in DB",
			"error", err,
			"post_id", postID,
		)
		utils.JSONError(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Post deleted successfully",
		"post_id", postID,
		"user_id", currentUserID,
	)
	utils.JSONSuccess(w, utils.SuccessResponse{Message: "post deleted"}, http.StatusOK)
}
