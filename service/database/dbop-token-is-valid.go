package database

// Given a username and a bearer auth token, checks if it matches token in DB
func (db *appdbimpl) TokenIsValid(username string, token string) (bool, error) {
	var expectedID string
	err := db.c.QueryRow("SELECT bearerAuthID from Users WHERE username=?", username).Scan(&expectedID)
	if err != nil{
		return false, err
	}
	if token == expectedID {
		return true, nil
	}
	return false, nil
}