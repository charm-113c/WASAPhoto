package database

func (db *appdbimpl) FollowUser(user1ID string, user2ID string) error {
	_, err := db.c.Exec("INSERT INTO Following (userID, followedUserID) VALUES (?, ?)", user1ID, user2ID)
	return err
}
