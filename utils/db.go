package utils

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func CloseDB(rows *sqlx.Rows) error {
	err := rows.Close()
	if err != nil {
		return err
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return err
}

func CommitOrRollback(tx *sqlx.Tx, message string, err error) {
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			LogError("CommitOrRollback", message, errRollback)
			return
		}
	} else {
		tx.Commit()
	}
}

func GenerateUUID() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		LogError("GenerateUUID", "Failed to generate UUID", err)
		// panic(err)
	}
	return id
}
