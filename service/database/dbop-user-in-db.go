package database

import (
	"database/sql"
	"errors"
)

// returns whether a user is present in the DB
func (db *appdbimpl) UserInDB(username string) (bool, error) {
	var temp string
	err := db.c.QueryRow("SELECT username FROM Users WHERE username=?", username).Scan(&temp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
