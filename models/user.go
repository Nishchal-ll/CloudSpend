package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Nishchal-ll/CloudSpend/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email" binding:"required,email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

// HashPassword hashes plaintext password with bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares password with bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthenticateUser verifies email and password
func AuthenticateUser(email, password string) (*User, error) {
	var user User
	var createdAtStr string
	query := `SELECT id, email, password_hash, name, created_at FROM users WHERE LOWER(email) = LOWER(?) LIMIT 1`
	row := database.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &createdAtStr)
	if err == sql.ErrNoRows {
		return nil, errors.New("invalid email or password")
	}
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// GetUserByEmail finds a user record
func GetUserByEmail(email string) (*User, error) {
	var user User
	var createdAtStr string
	query := `SELECT id, email, password_hash, name, created_at FROM users WHERE LOWER(email) = LOWER(?) LIMIT 1`
	row := database.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &createdAtStr)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
