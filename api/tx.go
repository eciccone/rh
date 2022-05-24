package repo

import (
	"database/sql"
	"errors"
	"fmt"
)

type transaction func(tx *sql.Tx) error

func Tx(db *sql.DB, fn transaction) error {
	if db == nil {
		return errors.New("repo.Tx() db is nil")
	}

	// begin tx
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("repo.Tx() unable to begin transaction: %v", err)
	}

	// execute tx, if error occurs then rollback
	if err = fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("repo.Tx() transaction failed: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repo.Tx() failed to commit transaction: %v", err)
	}

	return nil
}
