package database

import "time"

// defines user struct. The personal profile page for the user shows: the user’s photos (in reverse
// chronological order), how many photos have been uploaded, and the user’s followers and following.

type UserProfile struct {
	Username string
	Photos []Photo
	Nphotos int
	Followers []string
	Following []string
}

type Photo struct {
	PhotoID int  
	// UploaderID string
	BinaryData []byte
	Description string
	Likes int
	UploadDate time.Time 
	FileExtension string
	Comments int
}

// sorry for the disturbance
// ################### CREATE OPID GETPHOTO OR GETCOMMENTS
type Comment struct {
	CommentID int
	CommentText string
	UploaderID string
	UploadDate time.Time
}

