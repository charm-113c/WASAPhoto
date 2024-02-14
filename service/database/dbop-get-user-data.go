package database

// Find user by name
func (db *appdbimpl) GetUserData(username string) (UserData, error) {
	var userData UserData
	err := db.c.QueryRow("SELECT * FROM Users WHERE username=?", username).Scan(&userData.Username, &userData.bearerAuthID, &userData.UserID, &userData.Nphotos)
	return userData, err
}
