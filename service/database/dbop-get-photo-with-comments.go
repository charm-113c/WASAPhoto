package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPhotoWithComments(uploaderID string, photoID int) (PhotoWithComments, error) {
	// these being get operations, tx's aren't necessary
	var out PhotoWithComments
	// first get the photo
	pRow := db.c.QueryRow(`SELECT 
								p.photoID, u.username, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension
							FROM 
								Photos p LEFT JOIN Users AS u ON p.userID = u.userID
							WHERE 
								p.userID = ? AND p.photoID = ?`, uploaderID, photoID)
	// put photo in output struct
	err := pRow.Scan(&out.PhotoID, &out.Uploader, &out.BinaryData, &out.Description, &out.Likes, &out.UploadDate, &out.FileExtension)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// since an err is returned, out's values are irrelevant
			return out, err
		}
		// else no photos
		return PhotoWithComments{}, nil
		
	}
	// then its comments
	cRows, err := db.c.Query(`SELECT 
								c.content, c.commentID, u.username, c.uploadDate 
							FROM 
								Comments c LEFT JOIN Users u ON c.commenterID = u.userID
							WHERE 
								c.photoUploaderID = ? AND c.photoID = ?`, uploaderID, photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no comments shouldn't return an err
			return out, nil
		}
		return out, err
	}
	defer cRows.Close()
	// now do the same for comments
	for cRows.Next() {
		var comm Comment
		err = cRows.Scan(&comm.CommentText, &comm.CommentID, &comm.CommenterName, &comm.CommentDate)
		if err != nil {
			return out, err
		}
		// otherwise add comments to photo struct
		out.Comments = append(out.Comments, comm)
	}
	// uncaught loop error check
	if err = cRows.Err(); err != nil {
		return out, err
	}
	// and finally get liking users
	lRows, err := db.c.Query(`SELECT 
								u.username 
							FROM 
								Likes LEFT JOIN Users u ON Likes.likingUserID = u.userID 
							WHERE 
								Likes.uploaderID = ? AND Likes.photoID = ?`, uploaderID, photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no comments shouldn't return an err
			return out, nil
		}
		return out, err
	}
	defer lRows.Close()
	for lRows.Next() {
		var uname string
		err = lRows.Scan(&uname)
		if err != nil {
			return out, err
		}
		out.Likers = append(out.Likers, uname)
	}
	if err = lRows.Err(); err != nil {
		return out, err
	}

	return out, nil
}
