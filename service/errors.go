package service

import "errors"

// This file serves as a collection of all kinds of errors that package
// service can emit.

var (
	ErrInvalidEmailFormat    = errors.New("Invalid email format. Check RFC 5322.")
	ErrInvalidPasswordFormat = errors.New("Password format: Min. 8 chars, atleast 1 number" +
		"1 lowercase char, 1 uppercase char, 1 special char")
	ErrEmailExists          = errors.New("This email already exists.")
	ErrInvalidEmailPassword = errors.New("Invalid email/password combination.")
	ErrNilUser              = errors.New("Nil user.")
	ErrInvalidToken         = errors.New("Invalid token.")
	ErrNoInsert             = errors.New("Insertion failed, check logs.")
	ErrNoFetch              = errors.New("Fetching failed, check logs.")
	ErrNoRemove             = errors.New("Removal failed, check logs.")
)
