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
	statementInsertTodo          = "insert into todos (id, userId, task, status, deadline, completedAt) values ($1, $2, $3, $4, $5, $6);"
	statementSelectTodos         = "select * from todos where userId = $1 and status = 'pending' order by deadline;"
	statementUpdateTodo          = "update todos set task = $1, status = $2, deadline = $3, completedAt = $4 where id = $5 and userId = $6;"
	statementDeleteTodo          = "delete from todos where id = $1;"
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

// CreateTodo will create a new todo in the database with a
// new UUID.
func (db *Database) CreateTodo(t *entity.Todo) error {

	t.Id = uuid.New().String()
	err := db.execWithTransaction(statementInsertTodo, t.Id, t.UserId, t.Task, t.Status, t.Deadline, t.CompletedAt)
	return err
}

// GetAllTodos will return all todos for a given user using userId.
func (db *Database) GetAllPendingTodos(userId string) (*[]entity.Todo, error) {

	result := make([]entity.Todo, 0)
	err := db.queryWithTransaction(statementSelectTodos, func(rows *sql.Rows) error {
		for rows.Next() {
			var t entity.Todo
			err := rows.Scan(&t.Id, &t.UserId, &t.Task, &t.Status, &t.Deadline, &t.CompletedAt)
			if err != nil {
				return err
			}
			result = append(result, t)
		}
		return nil
	}, userId)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateTodo will update a todo in the database.
func (db *Database) UpdateTodo(t *entity.Todo) error {
	err := db.execWithTransaction(statementUpdateTodo, t.Task, t.Status, t.Deadline, t.CompletedAt, t.Id, t.UserId)
	return err
}

// DeleteTodo will delete a todo from the database.
func (db *Database) DeleteTodo(id string) error {
	return db.execWithTransaction(statementDeleteTodo, id)
}
