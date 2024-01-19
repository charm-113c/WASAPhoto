package database

import "github.com/gofrs/uuid"

func (db *appdbimpl) SetID(id uuid.UUID, username string) error {
	// called upon login, inserts given user's id in DB
	_, err := db.c.Exec("UPDATE Users SET userID = ? WHERE username = ?", id, username)
	return err
}