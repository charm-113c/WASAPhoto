package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPhotoData(userID string, photoID int) (Photo, error) {
	var out Photo
	row := db.c.QueryRow("SELECT FROM Photos (photoID, binaryData, description, likes, uploadDate, fileExtension, comments) WHERE userID = ? AND photoID = ?", userID, photoID)
	err := row.Scan(&out.PhotoID, &out.BinaryData, &out.Description, &out.Likes, &out.UploadDate, &out.FileExtension, &out.Comments)
	if err != nil {
		var empty Photo
		if errors.Is(err, sql.ErrNoRows) {
			return empty, nil
		}
		return empty, err
	}
	return out, nil
}
