package database

func (db *appdbimpl) GetPhotos(userID string) ([]Photo, error) {
	// get all of user's photos
	// rows, err := db.c.Query("SELECT p.photoID, p.postID, p.photoData, p.description, p.likes, p.uploadDate, p.fileExtension, COUNT(c.commentID) AS comments FROM Photos p LEFT JOIN Comments c ON p.photoID = c.photoID WHERE p.userID = ? GROUP BY p.photoID", userID)
	rows, err := db.c.Query("SELECT photoID, photoData, description, likes, uploadDate, fileExtension, comments FROM Photos WHERE userID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var photos []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoID, &photo.BinaryData, &photo.Description, &photo.Likes, &photo.UploadDate, &photo.FileExtension, &photo.Comments)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	err = rows.Err() // catch any possible error from Next(and thus incorrect/incomplete processing)
	if err != nil {
		return nil, err
	}

	return photos, nil
}
