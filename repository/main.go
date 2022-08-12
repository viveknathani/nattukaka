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

type Repository interface {
	userRepository
}
