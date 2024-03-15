package database

func (db *appdbimpl) UnbanUser(user1ID string, user2ID string) error {
	_, err := db.c.Exec("DELETE FROM Blacklist WHERE userID = ? AND bannedUserID = ?", user1ID, user2ID)
	return err
}
