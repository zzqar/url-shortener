package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"os"
	"url-shortener/internal/storage"
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

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const fn = "storage.sqlite.SaveURL"
	const sqlSaveURL = `INSERT INTO url (url, alias) VALUES (?,?)`

	stmt, err := s.db.Prepare(sqlSaveURL)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", fn, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s : %w", fn, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s : %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s : %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.sqlite.GetURL"
	const sqlGetURL = `SELECT url FROM url WHERE alias = ?`

	stmt, err := s.db.Prepare(sqlGetURL)
	if err != nil {
		return "", fmt.Errorf("%s : %w", fn, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		return "", fmt.Errorf("%s : %w", fn, err)
	}

	return resURL, nil
}
