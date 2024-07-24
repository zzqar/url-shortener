package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"
	const sqlInitPath = "db/migration/init.sql"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	dbInitSql, err := os.ReadFile(sqlInitPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	stmt, err := db.Prepare(string(dbInitSql))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	return &Storage{db: db}, nil

}
