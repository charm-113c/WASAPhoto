package database

func (db *appdbimpl) FollowingUser(user1ID string, user2 string) error {
	_, err := db.c.Exec("INSERT INTO Following (userID, followedUser) VALUES (?, ?)", user1ID, user2)
	return err
}