package database

// Insert row in Photos table with given input and Likes = 0
func (db *appdbimpl) UploadImage(photoID int, uID string, binData []byte, desc string, upDate string, ext string) error {
	_, err := db.c.Exec("INSERT INTO Photos (photoID, userID, photoData, description, likes, uploadDate, fileExtension) VALUES (?, ?, ?, ?, 0, ?, ?)", photoID, uID, binData, desc, upDate, ext)
	return err
}