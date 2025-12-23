package handlers

import "github.com/jackc/pgx/v5/pgxpool"

type Server struct {
	DB *pgxpool.Pool
}

// -- AUTH --
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// -- USER --
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// -- POST --
type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// -- COMMENT --
type CreateCommentInput struct {
	Content string `json:"content"`
}
