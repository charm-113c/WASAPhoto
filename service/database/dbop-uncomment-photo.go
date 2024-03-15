package database

func (db *appdbimpl) UncommentPhoto(uploaderID string, photoID int, commentID int) error {
	// composite operation: delete comment + decrement ncomments -> use transactions
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	// delete specified comment
	r, err := tx.Exec("DELETE FROM Comments WHERE photoUploaderID = ? AND photoID = ? and commentID = ?", uploaderID, photoID, commentID)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	// checking for idempotency
	n, err := r.RowsAffected()
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	if n == 0 {
		err = tx.Commit()
		return err
	}
	// decrement corresponding photo's ncomments in Photos
	_, err = tx.Exec("UPDATE Photos SET comments = comments - 1 WHERE userID = ? AND photoID = ?", uploaderID, photoID)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	// commit
	err = tx.Commit()
	return err
}
