package database

func (db *appdbimpl) UploadComment(comment string, commentID int, commenterID string, photoID int, uploaderID string, uploadDate string) error {
	// comopsite operation: using transactions
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// insert comment 
	_, err = tx.Exec("INSERT INTO Comments (content, commentID, photoID, photoUploaderID, uploadDate) VALUES (?,?,?,?,?,?)", comment, commentID, commenterID, photoID, uploaderID, uploadDate)
	if err != nil {
		return err
	}
	// increment photo's comment count by 1
	// done her for integrity, we have to be able to rollback
	_, err = tx.Exec("UPDATE Photos SET comments = comments + 1 WHERE photoID = ? AND userID = ?", photoID, uploaderID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
