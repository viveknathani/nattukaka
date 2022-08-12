package database

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/viveknathani/nattukaka/entity"
)

const dsn = "postgres://viveknathani:root@localhost:5432/nattukaka?sslmode=disable"

var db *Database

func TestMain(t *testing.M) {

	db = &Database{}
	err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// create tables
	_, err = db.pool.Exec("create table if not exists users(id uuid primary key,name varchar not null,email varchar(319) not null,password bytea not null);")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.pool.Exec("create table if not exists logs(id uuid primary key,userId uuid references users(id),latitude double precision,longitude double precision,activity varchar not null,startTime bigint not null,endTime bigint not null,notes varchar);")
	code := t.Run()
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
	os.Exit(code)
}

func TestCreateAndGetUser(t *testing.T) {

	u := &entity.User{
		Name:     "john",
		Email:    "john@gmail.com",
		Password: []byte("someHashedPwd44555"),
	}

	err := db.CreateUser(u)

	if err != nil {
		log.Fatal(err)
	}

	user, err := db.GetUser(u.Email)
	if err != nil {
		log.Fatal(err)
	}

	if user == nil {
		log.Fatal("Got nothing")
	}

	// since entity.User contains a byte array,
	// we cannot use the equality operator to test
	if !reflect.DeepEqual(*u, *user) {
		log.Println("Inequality.")
		log.Println("Created: ", u)
		log.Println("Got: ", user)
		log.Fatal()
	}

	// clean up
	err = db.DeleteUser(u.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func TestTodos(t *testing.T) {

	u := &entity.User{
		Name:     "alice",
		Email:    "alice@gmail.com",
		Password: []byte("someHashedPwd44555"),
	}

	err := db.CreateUser(u)

	if err != nil {
		log.Fatal(err)
	}

	deadline := "2022-08-13"
	td := &entity.Todo{
		UserId:      u.Id,
		Task:        "buy milk",
		Status:      "pending",
		Deadline:    &deadline,
		CompletedAt: nil,
	}

	err = db.CreateTodo(td)
	if err != nil {
		log.Fatal(err)
	}

	// update task
	td.Task = "buy vegetables"
	err = db.UpdateTodo(td)
	if err != nil {
		log.Fatal(err)
	}

	// update status and completedAt
	td.Status = "done"
	completed := "2022-08-14"
	td.CompletedAt = &completed
	err = db.UpdateTodo(td)
	if err != nil {
		log.Fatal(err)
	}

	err = db.DeleteTodo(td.Id)
	if err != nil {
		log.Fatal(err)
	}

	// clean up
	err = db.DeleteUser(u.Id)
	if err != nil {
		log.Fatal(err)
	}

}
