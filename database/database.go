package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

// Database contains the connection pool of SQL.
// Server should call Initialize before usage.
// Thanks to Go maintainers, the concurrency support is inbuilt
// so we do not need to manage connections on our own.
type Database struct {
	pool *sql.DB
}

// Initialize will open a database at given dataSourceName.
func (db *Database) Initialize(dataSourceName string) error {

	pool, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.pool = pool
	}
	return err
}

// Close is meant to free up resources, to be called when the
// server wants to shut down.
func (db *Database) Close() error {
	return db.pool.Close()
}

// queryWithTransaction runs the given query within a transaction.
// The caller can pass a function and form a closure to capture the values
// from row scanning. The caller should not worry about closing the rows.
func (db *Database) queryWithTransaction(prepared string, scanRows func(rows *sql.Rows) error, args ...interface{}) error {

	statement, err := db.pool.Prepare(prepared)
	if err != nil {
		return err
	}
	defer statement.Close()

	transaction, err := db.pool.Begin()
	if err != nil {
		return err
	}

	rows, err := transaction.Stmt(statement).Query(args...)
	if err != nil {
		rollError := transaction.Rollback()
		if rollError != nil {
			return errors.New(err.Error() + " " + rollError.Error())
		}
		return err
	}

	err = scanRows(rows)
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}
	return nil
}

// execWithTransaction will execute a given statement within a transaction.
func (db *Database) execWithTransaction(prepared string, args ...interface{}) error {

	statement, err := db.pool.Prepare(prepared)
	if err != nil {
		return err
	}

	transaction, err := db.pool.Begin()
	if err != nil {
		return err
	}

	_, err = transaction.Stmt(statement).Exec(args...)
	if err != nil {
		rollError := transaction.Rollback()
		if rollError != nil {
			return errors.New(err.Error() + " " + rollError.Error())
		}
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}
