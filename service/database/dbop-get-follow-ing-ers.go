package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetFollowing(userID string) ([]string, error) {
	// returns list of usernames of user's following (unsorted)
	var followings []string
	rows, err := db.c.Query("SELECT u.username FROM Following f LEFT JOIN Users u ON f.followedUserID = u.userID WHERE f.userID = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	// loop through each row
	for rows.Next() {
		var following string
		// put the value into following
		err = rows.Scan(&following)
		if err != nil {
			return nil, err
		}
		// update followings
		followings = append(followings, following)
	}
	// check for errors after loop
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return followings, nil
}

func (db *appdbimpl) GetFollowers(userID string) ([]string, error) {
	var followers []string
	rows, err := db.c.Query("SELECT u.username FROM Following f LEFT JOIN Users u on f.userID = u.userID WHERE f.followedUserID = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower string
		err = rows.Scan(&follower)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}
