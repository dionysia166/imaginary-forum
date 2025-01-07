package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Message holds data about a single message in a Thread.
type Message struct {
	ID      int
	Body    string
	Author  User
	DateAdded time.Time
}

// MessageModel holds a database handle for manipulating messages.
type MessageModel struct {
	DB *sql.DB
}

// NewMessageModel creates a Message table and returns a new MessageModel.
func NewMessageModel(db *sql.DB) (*MessageModel, error) {
	m := MessageModel{db}
	err := m.createTable()
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	return &m, nil
}

// createTable creates a Message table.
func (m *MessageModel) createTable() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS Messages (
		    id INTEGER PRIMARY KEY,
		    body TEXT NOT NULL,
		    author_id INTEGER NOT NULL REFERENCES Users,
		    thread_id INTEGER NOT NULL REFERENCES Threads,
		    created DATE NOT NULL
		);
	`
	_, err := m.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("creating Message table: %w", err)
	}
	return nil
}

// Insert inserts a new message in the Message table.
func (m *MessageModel) InsertMessage(
	body string, 
	threadId, 
	authorId int,
) (int, error) {
	stmt := `
		INSERT INTO Messages (body, thread_id, author_id, date_added)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`
	result, err := m.DB.Exec(stmt, body, threadId, authorId)
	if err != nil {
		return 0, fmt.Errorf("inserting new message in db: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert id: %w", err)
	}
	return int(id), nil
}
