/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserData(username string) (UserData, error)
	AddUser(name string, userid string) error
	SetNewName(currName string, newName string) error
	UserInDB(username string) (bool, error)
	UploadImage(photoID int, uID string, binData []byte, desc string, upDate string, ext string) error
	FollowUser(user1ID string, user2ID string) error
	HasBanned(user1ID string, user2ID string) (bool, error)
	GetPhotos(userID string) ([]Photo, error)
	GetFollowers(userID string) ([]string, error)
	GetFollowing(userID string) ([]string, error)
	GetFollowedPhotos(userID string) ([]Photo, error)
	BanUser(user1ID string, user2ID string) error
	UnfollowUser(user1ID string, user2ID string) error
	LikePhoto(uploaderID string, photoID int, likingUserID string) error
	GetPhotoData(userID string, photoID int) (Photo, error)
	UploadComment(comment string, commentID int, commenterID string, photoID int, uploaderID string, uploadDate string) error
	DeletePhoto(userID string, photoID int) error
	UnbanUser(user1ID string, user2ID string) error
	UnlikePhoto(uploaderID string, photoID int, likingUserID string) error
	UncommentPhoto(uploaderID string, photoID int, commentID int) error
	GetPhotoWithComments(uploaderID string, photoID int) (PhotoWithComments, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

type UserData struct {
	Username string
	UserID   string
	Nphotos  int
	TotNphotos int // also counts deleted photos
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// Note: as usernames can be changed, static userIDs are set to make updating the DB easier
		sqlStmt := `CREATE TABLE Users (
				username TEXT NOT NULL,
				userID TEXT NOT NULL PRIMARY KEY, 
				nphotos INTEGER,
				totnphotos INTEGER
			);
			
			CREATE TABLE Following (
				userID TEXT NOT NULL,
				followedUserID TEXT NOT NULL,
				PRIMARY KEY (userID, followedUserID),
				FOREIGN KEY (userID) REFERENCES Users(userID),
				FOREIGN KEY (followedUserID) REFERENCES Users(userID)
			);
			
			CREATE TABLE Blacklist (
				userID TEXT NOT NULL,
				bannedUserID TEXT NOT NULL,
				PRIMARY KEY (userID, bannedUserID),
				FOREIGN KEY (userID) REFERENCES Users(userID),
				FOREIGN KEY (bannedUserID) REFERENCES Users(userID)
			);
			
			CREATE TABLE Photos (
				photoID INTEGER NOT NULL,
				userID TEXT NOT NULL,
				photoData BLOB,
				description TEXT,
				likes INTEGER,
				uploadDate DATETIME,
				fileExtension TEXT,
				comments INTEGER,
				PRIMARY KEY (photoID, userID)
				FOREIGN KEY (userID) REFERENCES Users(userID)
			);
			
			CREATE TABLE Likes (
				uploaderID TEXT NOT NULL,
				photoID INTEGER NOT NULL,
				likingUserID TEXT NOT NULL,
				PRIMARY KEY (photoID, uploaderID, likingUserID),
				FOREIGN KEY (uploaderID) REFERENCES Users(userID),
				FOREIGN KEY (likingUserID) REFERENCES Users(userID),
				FOREIGN KEY (photoID) REFERENCES Photos(photoID)
			);
			
			CREATE TABLE Comments (
				content TEXT,
				commentID INTEGER NOT NULL,
				commenterID TEXT NOT NULL,
				photoID INTEGER NOT NULL,
				photoUploaderID TEXT NOT NULL,
				uploadDate DATETIME,
				PRIMARY KEY (commentID, photoID, photoUploaderID),
				FOREIGN KEY (photoID) REFERENCES Photos(photoID),
				FOREIGN KEY (commenterID) REFERENCES Users(userID),
				FOREIGN KEY (photoUploaderID) REFERENCES Users(userID)
			);`

		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

/*
In each photo, comments have a unique identifier: the number of comments+1 at the time
the comment is uploaded.
So, just like photoID, commentIDs aren't unique on their own, but require the photo
they belong to for uniqueness.
*/
