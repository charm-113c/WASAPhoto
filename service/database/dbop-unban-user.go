package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) UnbanUser(user1ID string, user2ID string) error {
	_, err := db.c.Exec("DELETE FROM Blacklist WHERE userID = ? AND bannedUserID = ?", user1ID, user2ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // for idempotency, ignore empty deletes
		return err
	}
	return nil
}
