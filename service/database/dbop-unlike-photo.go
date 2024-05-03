package database

func (db *appdbimpl) UnlikePhoto(uploaderID string, photoID int, likingUserID string) error {
	// composite op: deletes a row and decreases a value, so using transactions
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// delete Like triple
	r, err := tx.Exec("DELETE FROM Likes WHERE uploaderID = ? AND photoID = ? AND likingUserID = ?", uploaderID, photoID, likingUserID)
	if err != nil {
		// handle rollback first
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	// checking for idemptoency
	n, err := r.RowsAffected()
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	if n == 0 {
		// if no row found: do nothing
		err = tx.Commit()
		return err
	}
	// decrement likes in corresponding photo
	_, err = tx.Exec("UPDATE Photos SET  likes = likes - 1 WHERE userID = ? AND photoID = ?", uploaderID, photoID)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	// then commit transaction
	err = tx.Commit()
	return err
}
