package server

// This file defines how the incoming and outgoing JSON payloads look like.

type genericResponse struct {
	Message string `json:"message"`
}

type userSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
