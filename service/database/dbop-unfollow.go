package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) UnfollowUser(user1ID string, user2ID string) error {
	// Have user1 unfollow user2
	_, err := db.c.Exec("DELETE FROM Following WHERE user1ID = ? AND followedUserID = ?", user1ID, user2ID)
	if !errors.Is(err, sql.ErrNoRows) { // for idempotency
		return err
	}
	return nil
}