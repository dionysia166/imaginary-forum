package models

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User holds data about a user.
type User struct {
	ID       int
	Username string
	Email    string
	Password []byte
}

// UserModel holds a database handle for manipulating users.
type UserModel struct {
	DB *sql.DB
}

// NewUserModel creates a Threads table and returns a ThreadModel.
func NewUserModel(db *sql.DB) (*UserModel, error) {
	m := UserModel{db}
	err := m.createUserTable()
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	return &m, nil
}

// createUserTable creates a Users table.
func (m *UserModel) createUserTable() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY,
		    email TEXT NOT NULL,
		    username TEXT NOT NULL,
		    password TEXT NOT NULL,
		)
	`
	_, err := m.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("creating users table: %w", err)
	}
	return nil
}

// InsertUser inserts a new user into the database and returns its id.
func (m *UserModel) InsertUser(
	username string,
	email string,
	password string,
) (int, error) {
	emailExists, err := m.emailExists(email)
	if err != nil {
		return 0, fmt.Errorf("checking if email exists: %w", err)
	}
	if emailExists {
		return 0, ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, fmt.Errorf("generating hashed password: %w", err)
	}
	stmt := `
		INSERT INTO users (username, email, password) 
		VALUES (?, ?, ?)
	`
	result, err := m.DB.Exec(stmt, username, email, string(hashedPassword))
	if err != nil {
		return 0, fmt.Errorf("inserting row into database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert ID: %w", err)
	}
	return int(id), nil
}

// Authenticate checks if the email and password match a user in the database.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, password FROM users WHERE email = ?`

	result := m.DB.QueryRow(stmt, email)
	err := result.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, fmt.Errorf("querying database: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, fmt.Errorf("comparing password hashes: %w", err)
	}

	return id, nil
}

// emailExists checks if an email address is already in the database.
func (m *UserModel) emailExists(address string) (bool, error) {
	stmt := `SELECT id FROM users where email = ?`
	result := m.DB.QueryRow(stmt, address)
	var a string
	err := result.Scan(&a)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return true, fmt.Errorf("querying database: %w", err)
	}
	return true, nil
}

// GetUser will return a user based on id.
func (m *UserModel) GetUser(id int) (*User, error) {
	stmt := `
		SELECT id, username, email
		FROM users
		WHERE id = ?
	`
	row := m.DB.QueryRow(stmt, id)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}
	return &u, nil
}

