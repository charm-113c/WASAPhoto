package database

import "time"

// defines user struct. The personal profile page for the user shows: the user’s photos (in reverse
// chronological order), how many photos have been uploaded, and the user’s followers and following.

type UserProfile struct {
	Username  string
	Photos    []Photo
	Nphotos   int
	Followers []string
	Following []string
	Blacklist []string
}

type Photo struct {
	PhotoID       int
	Uploader      string
	BinaryData    []byte
	Description   string
	Likes         int
	Likers        string
	UploadDate    time.Time
	FileExtension string
	Comments      int
}

type Comment struct {
	CommentID     int
	CommentText   string
	CommenterName string
	CommentDate   time.Time
}

type PhotoWithComments struct {
	PhotoID       int
	Uploader      string
	BinaryData    []byte
	Description   string
	Likes         int
	Likers        []string
	UploadDate    time.Time
	FileExtension string
	Comments      []Comment
}
