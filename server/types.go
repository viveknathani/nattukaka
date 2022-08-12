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

type todoCreateRequest struct {
	Task     string `json:"task"`
	Deadline string `json:"deadline"`
}

type todoUpdateRequest struct {
	Id          string `json:"id"`
	Task        string `json:"task"`
	Status      string `json:"status"`
	Deadline    string `json:"deadline"`
	CompletedAt string `json:"completedAt"`
}

type todoDeleteRequest struct {
	Id string `json:"id"`
}
