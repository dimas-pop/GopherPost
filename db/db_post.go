package db

import (
	"gopher-post/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPostAll(dbpool *pgxpool.Pool, limit int, offset int) (*[]models.Post, error) {
	query := "SELECT id, title, content, user_id, created_at, updated_at FROM posts LIMIT $1 OFFSET $2"

	rows, err := dbpool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, err
}

func GetPostByID(dbpool *pgxpool.Pool, id string) (*models.Post, error) {
	query := "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = $1"

	var post models.Post
	err := dbpool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func CreatePostInDB(dbpool *pgxpool.Pool, title string, content string, user_id string) error {
	query := "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3)"

	_, err := dbpool.Exec(ctx, query, title, content, user_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePostByID(dbpool *pgxpool.Pool, title string, content string, id string) error {
	query := "UPDATE posts SET title = $1, content = $2 WHERE id = $3"

	_, err := dbpool.Exec(ctx, query, title, content, id)
	if err != nil {
		return err
	}

	return nil
}

func DeletePostByID(dbpool *pgxpool.Pool, id string) error {
	query := "DELETE FROM posts WHERE id = $1"

	_, err := dbpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
