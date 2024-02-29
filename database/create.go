package database

import (
	"database/sql"
	"os"
)

func create(db *sql.DB) error {
	sql, err := os.ReadFile("./schema.sql")

	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(string(sql))
	if err != nil {
		return err
	}
	return nil
}
