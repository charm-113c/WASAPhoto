package database

import "strings"

func (db *appdbimpl) GetPhotos(userID string) ([]Photo, error) {
	// get all of user's photos
	rows, err := db.c.Query(`SELECT 
								p.photoID, u.username, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, p.comments, GROUP_CONCAT(likers.username)
							FROM
								Photos p LEFT JOIN Users u ON p.userID = u.userID
								LEFT JOIN Likes l ON l.uploaderID = p.userID
								LEFT JOIN Users likers ON l.likingUserID = likers.userID
							WHERE 
								p.userID = ?
							GROUP BY
								p.photoID, u.username, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, p.comments`, userID)
	// while intimidating, this query is simple:
	// the first step is to get the photos belonging to user (first left join)
	// but we also need the list of users who like the photo, hence the next two left joins
	// group by + group concat is a combination that concatenates the names of the liking users into a single value (like concatenating strings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoID, &photo.Uploader, &photo.BinaryData, &photo.Description, &photo.Likes, &photo.UploadDate, &photo.FileExtension, &photo.Comments, &photo.Likers)
		if err != nil {
			if !strings.Contains(err.Error(),`"GROUP_CONCAT(likers.username)": converting NULL to string is unsupported`) {
				return nil, err
			}
			photo.Likers = ""
		}
		photos = append(photos, photo)
	}
	err = rows.Err() // catch any possible error from Next(and thus incorrect/incomplete processing)
	if err != nil {
		return nil, err
	}

	// we also need, for each photo, the list of users liking it
	// lRows, err := db.c.Query(`SELECT u.username FROM `) 

	return photos, nil
}
