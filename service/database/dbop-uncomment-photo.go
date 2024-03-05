package database

func (db *appdbimpl) UncommentPhoto(uploaderID string, photoID int, commentID int) error {
	// composite operation: delete comment + decrement ncomments -> use transactions
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	// delete specified comment
	_, err = tx.Exec("DELETE FROM Comments WHERE uploaderID = ? AND photoID = ? and commentID = ?", uploaderID, photoID, commentID)
	if err != nil {
		return err
	}

	// decrement corresponding photo's ncomments in Photos
	_, err = tx.Exec("UPDATE Photos SET comments = comments - 1 WHERE userID = ? AND photoID = ?", uploaderID, photoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}