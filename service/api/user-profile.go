package api

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
	BinaryData []byte
	Description string
	Likes int
	Comments []Comment 
	UploadDate time.Time 
	FileName string
	FileExtension string
}

type Comment struct {
	CommentID int
	CommentText string
	Uploader string
	UploadDate time.Time
}
