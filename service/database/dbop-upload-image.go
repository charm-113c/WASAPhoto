package database

// Insert row in Photos table with given input and Likes = 0
func (db *appdbimpl) UploadImage(photoID int, uID string, binData []byte, desc string, upDate string, ext string) error {
	// we upload photo and increase user's nphotos
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// upload img
	_, err = tx.Exec("INSERT INTO Photos (photoID, userID, photoData, description, likes, uploadDate, fileExtension, comments) VALUES (?, ?, ?, ?, 0, ?, ?, 0)", photoID, uID, binData, desc, upDate, ext)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	// increment nphotos
	_, err = tx.Exec("UPDATE Users SET nphotos = nphotos + 1 WHERE userID = ?", uID)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}
	err = tx.Commit()
	return err
}
