package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) UnlikePhoto(uploaderID string, photoID int, likingUserID string) error {
	// composite op: deletes a row and decreases a value, so using transactions
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// delete Like triple
	_, err = tx.Exec("DELETE FROM Likes WHERE uploaderID = ? AND photoID = ? AND likingUserID = ?", uploaderID, photoID, likingUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// idempotency at work
			return nil
		}
		return err
	}
	// decrement likes in corresponding photo
	_, err = tx.Exec("UPDATE Photos SET  likes = likes - 1 WHERE userID = ? AND photoID = ?", uploaderID, photoID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// then commit transaction
	err = tx.Commit()
	return err
}