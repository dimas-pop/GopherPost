package handlers

import (
	"encoding/json"
	"fmt"
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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Selamat Datang di GopherPost")
}

func (s *Server) GetPostAllHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))

	offSet := (limit * page) - limit

	posts, err := db.GetPostAll(s.DB, limit, offSet)
	if err != nil {
		http.Error(w, "Failed get post", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&posts)
}

func (s *Server) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	post, err := db.GetPostByID(s.DB, id)
	if err != nil {
		fmt.Println("ERROR DATABASE:", err)
		http.Error(w, "Failed get post", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&post)
}

func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if userID == "" {
		http.Error(w, "Missing userIDKey header", http.StatusUnauthorized)
		return
	}

	err = db.CreatePostInDB(s.DB, newPost.Title, newPost.Content, userID)
	if err != nil {
		http.Error(w, "failed create post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post Created"})
}

func (s *Server) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = db.UpdatePostByID(s.DB, input.Title, input.Content, id)
	if err != nil {
		http.Error(w, "Failed update post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated"})
}

func (s *Server) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.DeletePostByID(s.DB, id)
	if err != nil {
		http.Error(w, "Failed delete post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
}
