package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetFollowedPhotos(userID string) ([]Photo, error) {
	// join Following and Photo table on Following.followedUserID = Photo.userID WHERE Following.userID = userID
	rows, err := db.c.Query("SELECT p.photoID, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, p.comments FROM Photos p JOIN Following f ON p.userID = f.followedUserID WHERE f.userID = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var out []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoID, &photo.BinaryData, &photo.Description, &photo.Likes, &photo.UploadDate, &photo.FileExtension, &photo.Comments)
		if err != nil {
			return nil, err
		}
		out = append(out, photo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return out, nil
}