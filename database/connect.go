package database

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	name := "./scans.db"
	
	exists, err := exist(name)
	if err != nil {
		return nil, err
	}

	if !exists {
		db, err := connect(name)
		if err != nil {
			return nil, err
		}
		defer db.Close()

		err = create(db)
		if err != nil {
			return nil, err
		}
	}

	db, err := connect(name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connect(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=wal")
	if err != nil {
		return nil, err
	}

	return db, err
}

func exist(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
