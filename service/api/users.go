package api

import "time"

// Definition of Users

type Users struct {
	Username  string
	UserID    int
	Nphotos   int
	Photos    []Photo
	Following []string
	Followers []string
}

type Photo struct {
	PhotoData   []byte
	PhotoID     int
	Description string
	Likes       int
	Comments    []Comment
	UploadDate  time.Time
}

type Comment struct{
	Text string
	Uploader string
	UploadDate time.Time
}