package repository

import "github.com/viveknathani/nattukaka/entity"

// userRepository is the method set to operate on the User entity.
type userRepository interface {

	// CreateUser will create a new user in the database and will
	// have a newly generated UUID.
	CreateUser(u *entity.User) error

	// GetUser will fetch the first found user from the database.
	GetUser(email string) (*entity.User, error)

	// DeleteUser will delete a user specified by userId.
	DeleteUser(id string) error
}

type todoRepository interface {

	// CreateTodo will create a new todo in the database with a
	// new UUID.
	CreateTodo(t *entity.Todo) error

	// GetAllTodos will return all todos for a given user using userId.
	GetAllPendingTodos(userId string) (*[]entity.Todo, error)

	// UpdateTodo will update a todo in the database.
	UpdateTodo(t *entity.Todo) error

	// DeleteTodo will delete a todo from the database.
	DeleteTodo(id string) error
}

type noteRepository interface {
	CreateNote(n *entity.Note) error
	UpdateNote(n *entity.Note) error
	GetAllNotes(userId string) (*[]entity.Note, error)
	GetNote(id string, userId string) (*[]entity.Note, error)
}

type Repository interface {
	todoRepository
	userRepository
	noteRepository
}
