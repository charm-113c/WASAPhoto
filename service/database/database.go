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
	SetID(bearid string, username string) error
	TokenIsValid(username string, token string) (bool, error)
	SetNewName(currName string, newName string) error
	UserInDB(username string) (bool, error)
	UploadImage(photoID int, uID string, binData []byte, desc string, upDate string, ext string) error
	UpdateNphotos(username string, increase bool) error
	FollowingUser(user1ID string, user2 string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

type UserData struct {
	Username string
	bearerAuthID string
	UserID string
	Nphotos int
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
				username TEXT NOT NULL PRIMARY KEY,
				bearerAuthID TEXT,
				userID TEXT, 
				nphotos INTEGER
			);
			
			CREATE TABLE Following (
				userID TEXT,
				followedUser TEXT,
				PRIMARY KEY (userID, followedUser),
				FOREIGN KEY (userID) REFERENCES Users(userID),
				FOREIGN KEY (followedUser) REFERENCES Users(username)
			);
			
			CREATE TABLE Blacklist (
				userID TEXT,
				bannedUser TEXT,
				PRIMARY KEY (userID, bannedUser),
				FOREIGN KEY (userID) REFERENCES Users(userID),
				FOREIGN KEY (bannedUser) REFERENCES Users(username)
			);
			
			CREATE TABLE Photos (
				photoID INTEGER NOT NULL,
				userID TEXT NOT NULL,
				photoData BLOB,
				description TEXT,
				likes INTEGER,
				uploadDate DATETIME,
				fileExtension TEXT,
				PRIMARY KEY (photoID, userID)
				FOREIGN KEY (userID) REFERENCES Users(userID)
			);
			
			CREATE TABLE Likes (
				photoID INTEGER,
				userID TEXT,
				PRIMARY KEY (photoID, userID),
				FOREIGN KEY (userID) REFERENCES Users(userID),
				FOREIGN KEY (photoID) REFERENCES Photos(photoID)
			);
			
			CREATE TABLE Comments (
				commentID INTEGER NOT NULL,
				photoID INTEGER NOT NULL,
				userID TEXT NOT NULL,
				uploadDate DATETIME,
				PRIMARY KEY (commentID, photoID, userID)
				FOREIGN KEY (photoID) REFERENCES Photos(photoID),
				FOREIGN KEY (userID) REFERENCES Users(userID)
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
