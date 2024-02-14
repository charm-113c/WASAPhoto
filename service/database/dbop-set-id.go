package database

// called upon login, inserts given user's id in DB
func (db *appdbimpl) SetID(bearid string, username string) error {
	_, err := db.c.Exec("UPDATE Users SET bearerAuthID = ? WHERE username = ?", bearid, username)
	return err
}