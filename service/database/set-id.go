package database

func (db *appdbimpl) SetID(id int, username string) error {
	// called upon login, inserts given user's id in DB
	_, err := db.c.Exec("UPDATE Users SET userID = ? WHERE username = ?", id, username)
	return err
}