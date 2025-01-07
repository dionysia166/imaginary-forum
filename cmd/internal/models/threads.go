package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Thread holds data about a thread.
type Thread struct {
	ID        int
	Title     string
	Author    *User
	DateAdded time.Time
	Messages  []*Message
}

// ThreadModel holds a database handle to manipulate a Thread.
type ThreadModel struct {
	DB *sql.DB
}

// NewThreadModel creates a Threads table and returns a ThreadModel.
func NewThreadModel(db *sql.DB) (*ThreadModel, error) {
	m := ThreadModel{db}
	err := m.createTable()
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	return &m, nil
}

// createTable creates a Threads table.
func (m *ThreadModel) createTable() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS threads (
		    id INTEGER PRIMARY KEY,
		    title TEXT NOT NULL,
		    author_id INTEGER NOT NULL REFERENCES users,
		    created DATE NOT NULL
		)
	`
	_, err := m.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("creating threads table: %w", err)
	}
	return nil
}

// Insert inserts a new thread in the database and returns its id.
func (m *ThreadModel) Insert(title string, authorId int) (int, error) {
	stmt := `
		INSERT INTO threads (title, author_id, date_added)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`
	result, err := m.DB.Exec(stmt, title, authorId)
	if err != nil {
		return 0, fmt.Errorf("inserting new thread in db: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last thread id: %w", err)
	}
	return int(id), nil
}

// Get retrieves the thread with the given id from the database.
func (m *ThreadModel) Get(id int) (*Thread, error) {
	stmt := `
		SELECT t.id, t.title, t.date_added, u.id, u.username, u.email
		FROM threads T, users u
		WHERE t.author_id = u.id AND t.id = ?
	`
	row := m.DB.QueryRow(stmt, id)
	t, err := m.newThread(row, "ASC")
	if err != nil {
		return nil, fmt.Errorf("creating new thread: %w", err)
	}
	return t, nil
}

// Latests retrieves the 10 latests threads from the database.
func (m *ThreadModel) Latests() ([]*Thread, error) {
	stmt := `
		SELECT t.id, t.title, t.date_added, u.id, u.username, u.email
		FROM threads t, users u
		WHERE t.author_id = u.id
		ORDER BY t.date_added DESC
		LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("getting latests threads: %w", err)
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		t, err := m.newThread(rows, "DESC")
		if err != nil {
			return nil, fmt.Errorf("creating thread: %w", err)
		}
		threads = append(threads, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows for latests threads: %w", err)
	}

	return threads, nil
}

// scanner implements the Scan function.
type scanner interface {
	Scan(dest ...any) error
}

// newThread creates a new Thread. It also creates a User to represent
// the Thread's author, and the Messages associated with that Thread.
func (m *ThreadModel) newThread(s scanner, messageOrder string) (*Thread, error) {
	var (
		t Thread
		u User
	)
	err := s.Scan(
		&t.ID, &t.Title, &t.DateAdded,
		&u.ID, &u.Username, &u.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}
	t.Author = &u
	t.Messages, err = m.getMessages(t.ID, messageOrder)
	if err != nil {
		return nil, fmt.Errorf("getting messages with thread id %v: %w", t.ID, err)
	}
	return &t, nil
}

// getMessages retrieves all Messages related to the Thread with the given threadID.
// The value of order must be "ASC" or "DESC".
func (m *ThreadModel) getMessages(threadID int, order string) ([]*Message, error) {
	stmt := fmt.Sprintf(
		`
			SELECT m.id, m.body, m.date_added, u.id, u.username, u.email
			FROM messages m, users u
			WHERE m.author_id = u.id AND m.thread_id = ?
			ORDER BY m.date_added %v
		`,
		order,
	)
	rows, err := m.DB.Query(stmt, threadID, order)
	if err != nil {
		return nil, fmt.Errorf("getting messages: %w", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var (
			m Message
			u User
		)
		err := rows.Scan(
			&m.ID, &m.Body, &m.DateAdded,
			&u.ID, &u.Username, &u.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning message row: %w", err)
		}
		m.Author = u
		messages = append(messages, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over message rows: %w", err)
	}

	return messages, nil
}
