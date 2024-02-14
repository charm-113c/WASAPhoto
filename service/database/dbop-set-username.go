package database

func (db *appdbimpl) SetNewName(currName string, newName string) error {
	_, err := db.c.Exec("UPDATE Users SET username=? WHERE username=?", newName, currName)
	return err
}