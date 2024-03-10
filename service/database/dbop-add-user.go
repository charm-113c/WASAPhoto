package database

// Used only when a new user logs in: insert them in the DB
func (db *appdbimpl) AddUser(name, uID string) error {
	_, err := db.c.Exec("INSERT INTO Users (username, userID, nphotos) VALUES (?, ?, 0)", name, uID)
	return err
}
