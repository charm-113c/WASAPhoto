package database

// increases a user's number of photos by 1. Decreases by 1 if increase is set to false
func (db *appdbimpl) UpdateNphotos(username string, increase bool) error {
	if increase {	
		_, err := db.c.Exec("UPDATE Users SET nphotos = nphotos + 1 WHERE username = ?", username)
		return err
	} else  {
		_, err := db.c.Exec("UPDATE Users SET Nphotos = Nphotos - 1 WHERE username = ?", username)
		return err
	}
}