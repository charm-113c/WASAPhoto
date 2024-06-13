package database

import (
	"database/sql"
	"errors"
	"strings"
)

func (db *appdbimpl) GetFollowedPhotos(userID string) ([]Photo, error) {
	// join Following and Photo table on Following.followedUserID = Photo.userID WHERE Following.userID = userID
	rows, err := db.c.Query(`SELECT 
								p.photoID, u.username, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, p.comments, GROUP_CONCAT(likers.username) 
							FROM 
								Photos p 
								JOIN Following f ON p.userID = f.followedUserID 
								LEFT JOIN Users u ON f.followedUserID = u.userID
								LEFT JOIN Likes l ON p.userID = l.uploaderID AND p.photoID = l.photoID
    							LEFT JOIN Users AS likers ON l.likingUserID = likers.userID 
							WHERE 
								f.userID = ?
							GROUP BY
								p.photoID, u.username, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, p.comments`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	// loop through rows to get the data
	var out []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoID, &photo.Uploader, &photo.BinaryData, &photo.Description, &photo.Likes, &photo.UploadDate, &photo.FileExtension, &photo.Comments, &photo.Likers)
		if err != nil {
			if !strings.Contains(err.Error(), `"GROUP_CONCAT(likers.username)": converting NULL to string is unsupported`) {
				return nil, err
			}
			// the above check is to see if Likers is NULL (empty), in which case we simply set likers to an empty str
			photo.Likers = ""
		}
		out = append(out, photo)
	}
	// check for errors after loop
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return out, nil
}
