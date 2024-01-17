package database

// Find user by name
func (db *appdbimpl) FindByName(username string) (string, error) {
	var user string
	err := db.c.QueryRow("SELECT name FROM Users WHERE username=?", username).Scan(&user)
	return user, err
}
