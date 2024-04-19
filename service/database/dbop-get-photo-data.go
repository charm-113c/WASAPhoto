package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPhotoData(userID string, photoID int) (Photo, error) {
	var out Photo
	// get data from DB
	row := db.c.QueryRow(`SELECT 
							p.photoID, u.username, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, p.comments
						  FROM
						  	Photos p LEFT JOIN Users u ON p.userID = u.userID
						  WHERE 
						  	p.userID = ? AND p.photoID = ?`, userID, photoID)
	// put it in a variable
	err := row.Scan(&out.PhotoID, &out.Uploader, &out.BinaryData, &out.Description, &out.Likes, &out.UploadDate, &out.FileExtension, &out.Comments)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Photo{}, nil
		}
		return Photo{}, err
	}
	return out, nil
}
