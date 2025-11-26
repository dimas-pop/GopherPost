package handlers

import (
	"encoding/json"
	"fmt"
	"gopher-post/db"
	"gopher-post/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	result, err := db.GetCommentByPostID(s.DB, postID)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Failed get comment", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID := vars["id"]
	userID := r.Context().Value(middleware.UserIDKey).(string)

	err = db.CreateCommentInDB(s.DB, input.Content, userID, postID)
	if err != nil {
		http.Error(w, "Failed create comment", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment created"})
}

func (s *Server) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]

	currentUserID := r.Context().Value(middleware.UserIDKey).(string)

	userID, err := db.GetCommentOwnerID(s.DB, commentID)
	if err != nil {
		http.Error(w, "Failed get user_id", http.StatusBadRequest)
		return
	}

	if currentUserID != userID {
		http.Error(w, "Invalid user_id", http.StatusInternalServerError)
		return
	}

	err = db.DeleteCommentByID(s.DB, commentID)
	if err != nil {
		http.Error(w, "Failed delete comment", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted"})
}
