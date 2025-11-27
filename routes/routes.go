package routes

import (
	"gopher-post/handlers"
	"gopher-post/middleware"

	"github.com/gorilla/mux"

	_ "gopher-post/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(srv *handlers.Server) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", srv.LoginHandler).Methods("POST")
	router.HandleFunc("/register", srv.CreateUserHandler).Methods("POST")
	router.HandleFunc("/posts", srv.GetPostAllHandler).Methods("GET")
	router.HandleFunc("/posts/{id}", srv.GetPostByIDHandler).Methods("GET")
	router.HandleFunc("/posts/{id}/comments", srv.GetCommentHandler).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/posts", srv.CreatePostHandler).Methods("POST")
	api.HandleFunc("/posts/{id}", srv.UpdatePostHandler).Methods("PUT")
	api.HandleFunc("/posts/{id}", srv.DeletePostHandler).Methods("DELETE")
	api.HandleFunc("/posts/{id}/comments", srv.CreateCommentHandler).Methods("POST")
	api.HandleFunc("/comments/{id}", srv.DeleteCommentHandler).Methods("DELETE")

	api.HandleFunc("/users", srv.GetUserAllHandler).Methods("GET")
	api.HandleFunc("/users/{id}", srv.GetUserByIDHandler).Methods("GET")
	api.HandleFunc("/users/{id}", srv.UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id}", srv.DeleteUserHandler).Methods("DELETE")

	return router
}
