package db

import (
	"gopher-post/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCommentByPostID(dbpool *pgxpool.Pool, postID string) (*[]models.Comment, error) {
	query := "SELECT id, content, user_id, post_id, created_at FROM comments WHERE post_id = $1"

	rows, err := dbpool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return &comments, err
}

func GetCommentOwnerID(dbpool *pgxpool.Pool, id string) (string, error) {
	query := "SELECT user_id FROM comments WHERE id = $1"

	var userID string
	err := dbpool.QueryRow(ctx, query, id).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, err
}

func CreateCommentInDB(dbpool *pgxpool.Pool, comment string, userID string, postID string) error {
	query := "INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3)"

	_, err := dbpool.Exec(ctx, query, comment, userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCommentByID(dbpool *pgxpool.Pool, id string) error {
	query := "DELETE FROM comments WHERE id = $1"

	_, err := dbpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
