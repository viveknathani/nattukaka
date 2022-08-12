package database

// This file contains the implementation for the Repository interface.

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/viveknathani/nattukaka/entity"
)

const (
	statementInsertUser          = "insert into users (id, name, email, password) values ($1, $2, $3, $4);"
	statementSelectUserFromEmail = "select * from users where email = $1;"
	statementDeleteUser          = "delete from users where id = $1;"
)

// CreateUser will create a new user in the database and will
// have a newly generated UUID.
func (db *Database) CreateUser(u *entity.User) error {

	u.Id = uuid.New().String()
	err := db.execWithTransaction(statementInsertUser, u.Id, u.Name, u.Email, u.Password)
	return err
}

// GetUser will fetch the first found user from the database.
func (db *Database) GetUser(email string) (*entity.User, error) {

	var u entity.User
	exists := false
	err := db.queryWithTransaction(statementSelectUserFromEmail, func(rows *sql.Rows) error {

		//We iterate only once since we are interested in the first occurence
		if rows.Next() {
			err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password)
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

// DeleteUser will delete a user specified by userId.
func (db *Database) DeleteUser(id string) error {
	return db.execWithTransaction(statementDeleteUser, id)
}
