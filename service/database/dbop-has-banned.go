package database


// returns whether user1 has banned user2
func (db *appdbimpl) HasBanned(user1ID string, user2ID string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Blacklist WHERE userID = ? AND bannedUserID = ?", user1ID, user2ID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil  
	}
	return false, nil
}