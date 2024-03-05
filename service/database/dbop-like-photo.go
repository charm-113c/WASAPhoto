package database

func (db *appdbimpl) LikePhoto(uploaderID string, photoID int, likingUserID string) error {
	// composite operation: use transactions to allow rollback
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// insert Like tuple
	_, err = tx.Exec("INSERT INTO Likes (uploaderID, photoID, likingUserID) VALUES (?,?,?)", uploaderID, photoID, likingUserID)
	if err != nil {
		return err
	}
	// and increase number of comments of photo
	_, err = tx.Exec("UPDATE Photos SET likes = likes + 1 WHERE photoID = ? AND userID = ?", photoID, uploaderID) 
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
