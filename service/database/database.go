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
	FindByName(username string) (string, error)
	AddUser(name string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// create 3 tables: Users, Photos and Comments. Alright so it became more complex than I thought. Just a little more complex. Just a little.
		sqlStmt := `CREATE TABLE Users (
				username TEXT NOT NULL PRIMARY KEY, 
				userID INTEGER, 
				nphotos INTEGER
			);
			
			CREATE TABLE Following (
				username TEXT,
				followedUser TEXT,
				PRIMARY KEY (username, followedUser),
				FOREIGN KEY (username) REFERENCES Users(username),
				FOREIGN KEY (followedUser) REFERENCES Users(username)
			);
			
			CREATE TABLE Blacklist (
				username TEXT,
				bannedUser TEXT,
				PRIMARY KEY (username, bannedUser),
				FOREIGN KEY (username) REFERENCES Users(username),
				FOREIGN KEY (bannedUser) REFERENCES Users(username)
			);
			
			CREATE TABLE Photos (
				photoID INTEGER NOT NULL PRIMARY KEY,
				username TEXT,
				likes INTEGER,
				uploadDate DATETIME,
				description TEXT,
				FOREIGN KEY (username) REFERENCES Users(username)
			);
			
			CREATE TABLE Likes (
				photoID INTEGER,
				username TEXT,
				PRIMARY KEY (photoID, username),
				FOREIGN KEY (username) REFERENCES Users(username),
				FOREIGN KEY (photoID) REFERENCES Photos(photoID)
			);
			
			CREATE TABLE Comments (
				commentID INTEGER NOT NULL PRIMARY KEY,
				photoID INTEGER,
				username TEXT,
				uploadDate DATETIME,
				FOREIGN KEY (photoID) REFERENCES Photos(photoID),
				FOREIGN KEY (username) REFERENCES Users(username)
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
