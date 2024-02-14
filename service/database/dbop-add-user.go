package database

// Used only when a new user logs in: insert them in the DB
func (db *appdbimpl) AddUser(name, uID string) error {
	_, err := db.c.Exec("INSERT INTO Users (username, bearerAuthID, userID, nphotos) VALUES (?, '00000000-0000-0000-0000-000000000000', ?, 0)", name, uID)
	return err
}
