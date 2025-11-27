package handlers

import (
	"encoding/json"
	"gopher-post/db"
	"gopher-post/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

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
		JSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, &posts, http.StatusOK)
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
		JSONError(w, "Post not found or database error", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, &post, http.StatusOK)
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
		JSONError(w, "Bad request input", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		JSONError(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	err = db.CreatePostInDB(s.DB, newPost.Title, newPost.Content, userID)
	if err != nil {
		JSONError(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, SuccessResponse{Message: "post created"}, http.StatusCreated)
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
	id := vars["id"]

	var input UpdatePostInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		JSONError(w, "Bad request input", http.StatusBadRequest)
		return
	}

	err = db.UpdatePostByID(s.DB, input.Title, input.Content, id)
	if err != nil {
		JSONError(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, SuccessResponse{Message: "post updated"}, http.StatusOK)
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
	id := vars["id"]

	err := db.DeletePostByID(s.DB, id)
	if err != nil {
		JSONError(w, "Failed delete post", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, SuccessResponse{Message: "post deleted"}, http.StatusOK)
}
