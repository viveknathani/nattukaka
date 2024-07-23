package database

import (
	"database/sql"

	"github.com/viveknathani/nattukaka/types"
)

// SQL statements as constants
const (
	statementInsertUser        = `insert into users (name, email, public_id) values ($1, $2, $3) returning id`
	statementSelectUserByEmail = `select id, name, email, public_id from users where email = $1`
)

// InsertUser inserts a new user into the database
func (db *Database) InsertUser(u *types.User) error {
	err := db.execWithTransaction(statementInsertUser, u.Name, u.Email, u.PublicID)
	return err
}

// GetUserByEmail gets you a user by email
func (db *Database) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	exists := false
	err := db.query(statementSelectUserByEmail, func(rows *sql.Rows) error {
		// We iterate only once since we are interested in the first occurence
		if rows.Next() {
			err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PublicID)
			if err != nil {
				return err
			}
			exists = true
		}
		return nil
	}, email)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}
	return &u, nil
}
