package database

import (
	"database/sql"
	"embed"
)

var Embedded embed.FS
func create(db *sql.DB) error {
	sql, err := Embedded.ReadFile("schema.sql")

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
