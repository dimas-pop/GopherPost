package db

import (
	"gopher-post/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserAll(dbpool *pgxpool.Pool) (*[]models.User, error) {
	query := "SELECT id, name, email, created_at FROM users"

	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, err
}

func GetUserByID(dbpool *pgxpool.Pool, id string) (*models.User, error) {
	query := "SELECT id, name, email, password_hash, created_at FROM users WHERE id = $1"

	var user models.User
	err := dbpool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetUserByEmail(dbpool *pgxpool.Pool, email string) (*models.User, error) {
	query := "SELECT id, password_hash FROM users WHERE email = $1"

	var user models.User
	err := dbpool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func CheckEmailExists(dbpool *pgxpool.Pool, email string) (bool, error) {
	query := "SELECT id FROM users WHERE email = $1"

	var id string
	err := dbpool.QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CreateUserInDB(dbpool *pgxpool.Pool, name string, email string, pass_hash string) error {
	query := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)"

	_, err := dbpool.Exec(ctx, query, name, email, pass_hash)
	if err != nil {
		return err
	}

	return err
}

func UpdateUserByID(dbpool *pgxpool.Pool, name string, email string, id string) error {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"

	_, err := dbpool.Exec(ctx, query, name, email, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserByID(dbpool *pgxpool.Pool, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := dbpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
