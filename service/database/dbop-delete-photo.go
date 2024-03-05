package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) DeletePhoto(userID string, photoID int) error {
	// operations are complex, transactions (tx) allow us to rollback if something goes wrong
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	// delete from Photos table
	rows, err := tx.Exec("DELETE FROM Photos WHERE userID = ? AND photoID = ?", userID, photoID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err // no db-changing ops are performed prior -> no rollback necessary
	}
	nRows, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if nRows == 0 {
		// implementing idempotency: if photo not found, do nothing
		return nil
	}

	// else proceed to eliminate likes, comments
	_, err = tx.Exec("DELETE FROM Likes WHERE uploaderID = ? AND photoID = ?", userID, photoID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { tx.Rollback(); return err}
	_, err = tx.Exec("DELETE FROM Comments WHERE photoUploaderID = ? AND photoID = ?", userID, photoID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { tx.Rollback(); return err}
	// and decrement nphotos 
	_, err = tx.Exec("UPDATE Users SET nphotos = nphotos - 1 WHERE userID = ?", userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { tx.Rollback(); return err}

	// and commit transactions 
	err = tx.Commit()
	return err 
}