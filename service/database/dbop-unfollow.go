package database

func (db *appdbimpl) UnfollowUser(user1ID string, user2ID string) error {
	// Have user1 unfollow user2
	_, err := db.c.Exec("DELETE FROM Following WHERE userID = ? AND followedUserID = ?", user1ID, user2ID)
	return err
}
