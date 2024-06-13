package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetBlacklist(userID string) ([]string, error) {
	var blacklist []string
	rows, err := db.c.Query(`SELECT u.username 
								FROM Blacklist b LEFT JOIN Users u
								ON b.bannedUserID = u.userID
								WHERE b.userID = ?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var banned string
		if err = rows.Scan(&banned); err != nil {
			return nil, err
		}
		blacklist = append(blacklist, banned)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blacklist, nil
}
