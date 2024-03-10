package database

func (db *appdbimpl) BanUser(user1ID string, user2ID string) error {
	_, err := db.c.Exec("INSERT INTO Blacklist (userID, bannedUserID) VALUES (?, ?)", user1ID, user2ID)
	return err
}