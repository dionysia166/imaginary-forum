package models

import (
	"database/sql"
)

// NewModels creates all models necessary for the application.
func NewModels(db *sql.DB) (*ThreadModel, *UserModel, *MessageModel, error) {
	threadModel, err := NewThreadModel(db)
	if err != nil {
		return nil, nil, nil, err
	}
	userModel, err := NewUserModel(db)
	if err != nil {
		return nil, nil, nil, err
	}
	postModel, err := NewMessageModel(db)
	if err != nil {
		return nil, nil, nil, err
	}
	return threadModel, userModel, postModel, nil
}
